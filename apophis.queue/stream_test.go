package apophis

import (
	"context"
	"errors"
	"io"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ninesbr/sheeps.toolkit.go/apophis.queue/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type fakeSubscribeStream struct {
	ctx    context.Context
	sendFn func(*pb.SubscribeMessage) error
	recvFn func() (*pb.SubscribeMessage, error)
}

func (f *fakeSubscribeStream) Send(msg *pb.SubscribeMessage) error {
	if f.sendFn == nil {
		return nil
	}
	return f.sendFn(msg)
}

func (f *fakeSubscribeStream) Recv() (*pb.SubscribeMessage, error) {
	if f.recvFn == nil {
		return nil, io.EOF
	}
	return f.recvFn()
}

func (f *fakeSubscribeStream) Header() (metadata.MD, error) { return nil, nil }
func (f *fakeSubscribeStream) Trailer() metadata.MD         { return nil }
func (f *fakeSubscribeStream) CloseSend() error             { return nil }
func (f *fakeSubscribeStream) Context() context.Context     { return f.ctx }

func (f *fakeSubscribeStream) SendMsg(msg any) error {
	return f.Send(msg.(*pb.SubscribeMessage))
}

func (f *fakeSubscribeStream) RecvMsg(msg any) error {
	received, err := f.Recv()
	if err != nil {
		return err
	}
	proto.Merge(msg.(*pb.SubscribeMessage), received)
	return nil
}

type fakePubSubClient struct {
	pb.PubSubServiceClient
	pingFn      func(context.Context) error
	publishFn   func(context.Context, *pb.PubMessageRequest) (*pb.PubMessageResponse, error)
	subscribeFn func(context.Context) (grpc.BidiStreamingClient[pb.SubscribeMessage, pb.SubscribeMessage], error)
}

func (f *fakePubSubClient) Ping(ctx context.Context, _ *pb.PingRequest, _ ...grpc.CallOption) (*pb.PingResponse, error) {
	if f.pingFn != nil {
		if err := f.pingFn(ctx); err != nil {
			return nil, err
		}
	}
	return &pb.PingResponse{}, nil
}

func (f *fakePubSubClient) Publish(ctx context.Context, req *pb.PubMessageRequest, _ ...grpc.CallOption) (*pb.PubMessageResponse, error) {
	if f.publishFn == nil {
		return &pb.PubMessageResponse{}, nil
	}
	return f.publishFn(ctx, req)
}

func (f *fakePubSubClient) Subscribe(ctx context.Context, _ ...grpc.CallOption) (grpc.BidiStreamingClient[pb.SubscribeMessage, pb.SubscribeMessage], error) {
	return f.subscribeFn(ctx)
}

func TestDeliveryConfirmationSendsOnlyOnce(t *testing.T) {
	t.Parallel()

	var mu sync.Mutex
	var sent []*pb.SubscribeMessage
	stream := &fakeSubscribeStream{
		ctx: context.Background(),
		sendFn: func(msg *pb.SubscribeMessage) error {
			mu.Lock()
			defer mu.Unlock()
			sent = append(sent, proto.Clone(msg).(*pb.SubscribeMessage))
			return nil
		},
	}
	sender := &serializedStreamSender{stream: stream}
	original := &pb.SubscribeMessage{
		Id:      "message-1",
		Headers: map[string]string{"original": "value"},
	}
	confirmation := newDeliveryConfirmation(
		context.Background(),
		sender,
		original,
		make(chan error, 1),
	)

	var callers sync.WaitGroup
	for range 32 {
		callers.Add(1)
		go func() {
			defer callers.Done()
			confirmation.confirm(pb.MessageCommit_RETRY, map[string]string{"attempt": "2"})
		}()
	}
	callers.Wait()

	mu.Lock()
	defer mu.Unlock()
	if len(sent) != 1 {
		t.Fatalf("expected exactly one confirmation, got %d", len(sent))
	}
	if sent[0].Commit != pb.MessageCommit_RETRY {
		t.Fatalf("expected RETRY, got %s", sent[0].Commit)
	}
	if sent[0].Headers["attempt"] != "2" {
		t.Fatalf("expected retry header in confirmation")
	}
	if _, changed := original.Headers["attempt"]; changed {
		t.Fatalf("received message was mutated while building confirmation")
	}
}

func TestSerializedStreamSenderPreventsConcurrentSend(t *testing.T) {
	t.Parallel()

	var active atomic.Int32
	var concurrent atomic.Bool
	var sent atomic.Int32
	stream := &fakeSubscribeStream{
		ctx: context.Background(),
		sendFn: func(*pb.SubscribeMessage) error {
			if active.Add(1) != 1 {
				concurrent.Store(true)
			}
			time.Sleep(time.Millisecond)
			active.Add(-1)
			sent.Add(1)
			return nil
		},
	}
	sender := &serializedStreamSender{stream: stream}

	var callers sync.WaitGroup
	for range 32 {
		callers.Add(1)
		go func() {
			defer callers.Done()
			if err := sender.Send(&pb.SubscribeMessage{}); err != nil {
				t.Errorf("Send() error = %v", err)
			}
		}()
	}
	callers.Wait()

	if concurrent.Load() {
		t.Fatal("stream received concurrent Send calls")
	}
	if sent.Load() != 32 {
		t.Fatalf("expected 32 sends, got %d", sent.Load())
	}
}

func TestDeliveryConfirmationTimeoutDiscardsMessage(t *testing.T) {
	t.Parallel()

	sent := make(chan *pb.SubscribeMessage, 1)
	stream := &fakeSubscribeStream{
		ctx: context.Background(),
		sendFn: func(msg *pb.SubscribeMessage) error {
			sent <- proto.Clone(msg).(*pb.SubscribeMessage)
			return nil
		},
	}
	confirmation := newDeliveryConfirmation(
		context.Background(),
		&serializedStreamSender{stream: stream},
		&pb.SubscribeMessage{Id: "message-1"},
		make(chan error, 1),
	)
	confirmation.startTimeout(time.Millisecond)

	select {
	case msg := <-sent:
		if msg.Commit != pb.MessageCommit_DISCARD {
			t.Fatalf("expected DISCARD, got %s", msg.Commit)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout did not confirm the message")
	}
}

func TestWatchingCancellationEndsReceive(t *testing.T) {
	t.Parallel()

	receiveEnded := make(chan struct{})
	signed := make(chan struct{})
	var signOnce sync.Once
	client := &fakePubSubClient{}
	client.subscribeFn = func(ctx context.Context) (grpc.BidiStreamingClient[pb.SubscribeMessage, pb.SubscribeMessage], error) {
		return &fakeSubscribeStream{
			ctx: ctx,
			sendFn: func(msg *pb.SubscribeMessage) error {
				if msg.Sign != nil {
					signOnce.Do(func() { close(signed) })
				}
				return nil
			},
			recvFn: func() (*pb.SubscribeMessage, error) {
				<-ctx.Done()
				close(receiveEnded)
				return nil, ctx.Err()
			},
		}, nil
	}
	a := &apophis{
		ops:    NewOptions(WithQueueName("test")),
		client: client,
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- a.watching(ctx, make(chan *MessageResponse[any]))
	}()

	select {
	case <-signed:
	case <-time.After(time.Second):
		t.Fatal("subscription sign was not sent")
	}
	cancel()

	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context cancellation, got %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("watching did not stop after cancellation")
	}
	select {
	case <-receiveEnded:
	case <-time.After(time.Second):
		t.Fatal("Recv remained blocked after cancellation")
	}
}

func TestSubscribeReconnectsAfterEOF(t *testing.T) {
	t.Parallel()

	attempts := make(chan int, 2)
	var attempt atomic.Int32
	client := &fakePubSubClient{}
	client.subscribeFn = func(ctx context.Context) (grpc.BidiStreamingClient[pb.SubscribeMessage, pb.SubscribeMessage], error) {
		current := int(attempt.Add(1))
		attempts <- current
		return &fakeSubscribeStream{
			ctx: ctx,
			recvFn: func() (*pb.SubscribeMessage, error) {
				if current == 1 {
					return nil, io.EOF
				}
				<-ctx.Done()
				return nil, ctx.Err()
			},
		}, nil
	}
	a := &apophis{
		ops: NewOptions(
			WithQueueName("test"),
			WithReconnectInterval(time.Millisecond),
		),
		client: client,
	}

	messages, cancel := a.subscribe(context.Background())
	defer cancel()

	for expected := 1; expected <= 2; expected++ {
		select {
		case got := <-attempts:
			if got != expected {
				t.Fatalf("subscription attempt = %d, want %d", got, expected)
			}
		case <-time.After(time.Second):
			t.Fatalf("subscription attempt %d did not happen", expected)
		}
	}

	cancel()
	select {
	case _, open := <-messages:
		if open {
			t.Fatal("message channel remained open after cancellation")
		}
	case <-time.After(time.Second):
		t.Fatal("subscription did not stop after cancellation")
	}
}

func TestWaitForReconnectStopsOnCancellation(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan bool, 1)
	go func() {
		done <- waitForReconnect(ctx, time.Hour)
	}()
	cancel()

	select {
	case reconnect := <-done:
		if reconnect {
			t.Fatal("wait requested reconnect after context cancellation")
		}
	case <-time.After(time.Second):
		t.Fatal("reconnect wait did not react to cancellation")
	}
}

func TestShouldReconnect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "clean stream end", err: io.EOF, want: true},
		{name: "unavailable", err: status.Error(codes.Unavailable, "offline"), want: true},
		{name: "unknown", err: status.Error(codes.Unknown, "temporary"), want: true},
		{name: "permission denied", err: status.Error(codes.PermissionDenied, "denied"), want: false},
		{name: "invalid argument", err: status.Error(codes.InvalidArgument, "invalid"), want: false},
		{name: "canceled", err: context.Canceled, want: false},
		{name: "no error", err: nil, want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := shouldReconnect(test.err); got != test.want {
				t.Fatalf("shouldReconnect(%v) = %t, want %t", test.err, got, test.want)
			}
		})
	}
}
