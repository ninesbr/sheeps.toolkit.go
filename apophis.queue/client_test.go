package apophis

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ninesbr/sheeps.toolkit.go/apophis.queue/pb"
	"google.golang.org/grpc"
)

type fakeControlClient struct {
	pb.PubSubServiceClient
	call func(context.Context) error
}

func (f *fakeControlClient) Ping(ctx context.Context, _ *pb.PingRequest, _ ...grpc.CallOption) (*pb.PingResponse, error) {
	if err := f.call(ctx); err != nil {
		return nil, err
	}
	return &pb.PingResponse{}, nil
}

func (f *fakeControlClient) Create(ctx context.Context, _ *pb.PubRequest, _ ...grpc.CallOption) (*pb.PubResponse, error) {
	if err := f.call(ctx); err != nil {
		return nil, err
	}
	return &pb.PubResponse{}, nil
}

func (f *fakeControlClient) Drop(ctx context.Context, _ *pb.DropRequest, _ ...grpc.CallOption) (*pb.PubResponse, error) {
	if err := f.call(ctx); err != nil {
		return nil, err
	}
	return &pb.PubResponse{}, nil
}

func TestControlOperationsHonorRequestTimeout(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		call func(*apophis) error
	}{
		{name: "ping", call: func(a *apophis) error { return a.Ping() }},
		{name: "create", call: func(a *apophis) error { return a.Create() }},
		{name: "drop", call: func(a *apophis) error { return a.Drop(false) }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			const timeout = 100 * time.Millisecond
			client := &fakeControlClient{
				call: func(ctx context.Context) error {
					deadline, ok := ctx.Deadline()
					if !ok {
						t.Fatal("operation context has no deadline")
					}
					remaining := time.Until(deadline)
					if remaining <= 0 || remaining > timeout {
						t.Fatalf("deadline remaining = %s, want within (0, %s]", remaining, timeout)
					}

					<-ctx.Done()
					return ctx.Err()
				},
			}
			a := &apophis{
				ops: NewOptions(
					WithQueueName("queue-1"),
					WithRequestTimeout(timeout),
				),
				client: client,
			}

			if err := test.call(a); !errors.Is(err, context.DeadlineExceeded) {
				t.Fatalf("operation error = %v, want context deadline exceeded", err)
			}
		})
	}
}
