package apophis

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/ninesbr/sheeps.toolkit.go/apophis.queue/pb"
	"google.golang.org/protobuf/proto"
)

type recordingApophis struct {
	publishRequest *MessageRequest
	publishErr     error
	publishCalls   int
}

var _ ApophisInterface = (*recordingApophis)(nil)

func (f *recordingApophis) Ping() error { return nil }

func (f *recordingApophis) publish(req *MessageRequest) error {
	f.publishCalls++
	f.publishRequest = req
	return f.publishErr
}

func (f *recordingApophis) subscribe(context.Context) (<-chan *MessageResponse[any], context.CancelFunc) {
	return nil, func() {}
}

func (f *recordingApophis) Drop(bool) error { return nil }
func (f *recordingApophis) Create() error   { return nil }
func (f *recordingApophis) Close() error    { return nil }

func TestPublisherBuildsJSONRequest(t *testing.T) {
	t.Parallel()

	type message struct {
		Text string `json:"text"`
	}
	client := &recordingApophis{}
	publisher := NewPublisher[message](client)

	err := publisher.Publish(
		&message{Text: "hello"},
		WithContentTypeRequest("application/problem+json; charset=utf-8"),
		WithHeaderRequest("source", "test"),
		WithTagRequest("one", "two"),
		WithCustomIDRequest("custom-id"),
		WithTrackingIDRequest("tracking-id"),
		WithForceCreateRequest(true),
	)
	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}

	req := client.publishRequest
	if req == nil {
		t.Fatal("publish request was not sent")
	}
	if req.ContentType != "application/problem+json; charset=utf-8" {
		t.Fatalf("ContentType = %q", req.ContentType)
	}
	if req.Headers["source"] != "test" {
		t.Fatalf("source header = %q", req.Headers["source"])
	}
	if !reflect.DeepEqual(req.Tags, []string{"one", "two"}) {
		t.Fatalf("Tags = %v", req.Tags)
	}
	if req.CustomID != "custom-id" || req.TrackingID != "tracking-id" {
		t.Fatalf("IDs were not copied: custom=%q tracking=%q", req.CustomID, req.TrackingID)
	}
	if !req.ForceCreate {
		t.Fatal("ForceCreate was not copied")
	}

	var decoded message
	if err := json.Unmarshal(req.Body, &decoded); err != nil {
		t.Fatalf("published body is not JSON: %v", err)
	}
	if decoded.Text != "hello" {
		t.Fatalf("published message = %#v", decoded)
	}
}

func TestPublisherRejectsInvalidInput(t *testing.T) {
	t.Parallel()

	type message struct {
		Text string `json:"text"`
	}
	validMessage := &message{Text: "hello"}

	tests := []struct {
		name      string
		publisher *Publisher[message]
		msg       *message
		options   []func(*MessageRequestOptions)
		wantError string
	}{
		{
			name:      "nil publisher",
			publisher: nil,
			msg:       validMessage,
			wantError: "publisher client is nil",
		},
		{
			name:      "nil client",
			publisher: NewPublisher[message](nil),
			msg:       validMessage,
			wantError: "publisher client is nil",
		},
		{
			name:      "nil message",
			publisher: NewPublisher[message](&recordingApophis{}),
			wantError: "message is nil",
		},
		{
			name:      "nil option",
			publisher: NewPublisher[message](&recordingApophis{}),
			msg:       validMessage,
			options:   []func(*MessageRequestOptions){nil},
			wantError: "option 0 is nil",
		},
		{
			name:      "non JSON content type",
			publisher: NewPublisher[message](&recordingApophis{}),
			msg:       validMessage,
			options:   []func(*MessageRequestOptions){WithContentTypeRequest("application/protobuf")},
			wantError: "not compatible with JSON",
		},
		{
			name:      "invalid content type",
			publisher: NewPublisher[message](&recordingApophis{}),
			msg:       validMessage,
			options:   []func(*MessageRequestOptions){WithContentTypeRequest("not a media type")},
			wantError: "invalid content type",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := test.publisher.Publish(test.msg, test.options...)
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Fatalf("Publish() error = %v, want error containing %q", err, test.wantError)
			}
		})
	}
}

func TestPublisherWrapsMarshalAndClientErrors(t *testing.T) {
	t.Parallel()

	t.Run("marshal", func(t *testing.T) {
		t.Parallel()
		type invalidMessage struct {
			Channel chan int `json:"channel"`
		}
		err := NewPublisher[invalidMessage](&recordingApophis{}).Publish(&invalidMessage{})
		if err == nil || !strings.Contains(err.Error(), "marshal publish message") {
			t.Fatalf("Publish() error = %v", err)
		}
	})

	t.Run("client", func(t *testing.T) {
		t.Parallel()
		sentinel := errors.New("server unavailable")
		client := &recordingApophis{publishErr: sentinel}
		message := struct{ Text string }{Text: "hello"}
		err := NewPublisher[struct{ Text string }](client).Publish(&message)
		if !errors.Is(err, sentinel) {
			t.Fatalf("Publish() error = %v, want wrapped sentinel", err)
		}
	})
}

func TestClientPublishAppliesTimeoutAndForceCreate(t *testing.T) {
	t.Parallel()

	var captured *pb.PubMessageRequest
	var remaining time.Duration
	client := &fakePubSubClient{}
	client.publishFn = func(ctx context.Context, req *pb.PubMessageRequest) (*pb.PubMessageResponse, error) {
		deadline, ok := ctx.Deadline()
		if !ok {
			t.Fatal("Publish context has no deadline")
		}
		remaining = time.Until(deadline)
		captured = proto.Clone(req).(*pb.PubMessageRequest)
		return &pb.PubMessageResponse{MsgID: "message-1"}, nil
	}

	const timeout = 250 * time.Millisecond
	a := &apophis{
		ops: NewOptions(
			WithQueueName("queue-1"),
			WithQueueDurable(true),
			WithQueueKeepMessages(true),
			WithQueueTags("queue-tag"),
			WithPublishTimeout(timeout),
		),
		client: client,
	}
	err := a.publish(&MessageRequest{
		ContentType: "application/json",
		Body:        []byte(`{"text":"hello"}`),
		ForceCreate: true,
	})
	if err != nil {
		t.Fatalf("publish() error = %v", err)
	}
	if remaining <= 0 || remaining > timeout {
		t.Fatalf("deadline remaining = %s, want within (0, %s]", remaining, timeout)
	}
	if captured == nil || captured.ForceCreate == nil {
		t.Fatal("ForceCreate was not sent to gRPC")
	}
	if captured.ForceCreate.Uniqid != "queue-1" || !captured.ForceCreate.Durable || !captured.ForceCreate.KeepMessages {
		t.Fatalf("ForceCreate config = %#v", captured.ForceCreate)
	}
	if !reflect.DeepEqual(captured.ForceCreate.Tags, []string{"queue-tag"}) {
		t.Fatalf("ForceCreate tags = %v", captured.ForceCreate.Tags)
	}
}

func TestClientPublishDoesNotForceCreateByDefault(t *testing.T) {
	t.Parallel()

	var captured *pb.PubMessageRequest
	client := &fakePubSubClient{}
	client.publishFn = func(_ context.Context, req *pb.PubMessageRequest) (*pb.PubMessageResponse, error) {
		captured = proto.Clone(req).(*pb.PubMessageRequest)
		return &pb.PubMessageResponse{}, nil
	}
	a := &apophis{
		ops:    NewOptions(WithQueueName("queue-1")),
		client: client,
	}

	if err := a.publish(&MessageRequest{ContentType: "application/json", Body: []byte(`{}`)}); err != nil {
		t.Fatalf("publish() error = %v", err)
	}
	if captured == nil {
		t.Fatal("publish request was not sent")
	}
	if captured.ForceCreate != nil {
		t.Fatalf("ForceCreate = %#v, want nil without explicit option", captured.ForceCreate)
	}
}

func TestClientPublishHonorsTimeout(t *testing.T) {
	t.Parallel()

	client := &fakePubSubClient{}
	client.publishFn = func(ctx context.Context, _ *pb.PubMessageRequest) (*pb.PubMessageResponse, error) {
		<-ctx.Done()
		return nil, ctx.Err()
	}
	a := &apophis{
		ops: NewOptions(
			WithQueueName("queue-1"),
			WithPublishTimeout(10*time.Millisecond),
		),
		client: client,
	}

	err := a.publish(&MessageRequest{ContentType: "application/json", Body: []byte(`{}`)})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("publish() error = %v, want context deadline exceeded", err)
	}
}

func TestOptionsRejectInvalidPublishTimeout(t *testing.T) {
	t.Parallel()

	ops := NewOptions(
		WithHost("localhost"),
		WithPort(50051),
		WithQueueName("queue-1"),
		WithPublishTimeout(0),
	)
	if err := ops.Validate(); err == nil || !strings.Contains(err.Error(), "publish timeout") {
		t.Fatalf("Validate() error = %v", err)
	}
}
