package apophis

import (
	"context"
	"fmt"
	"time"

	"github.com/ninesbr/sheeps.toolkit.go/apophis/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type ApophisInterface interface {
	Ping() error
	Publish(ctx context.Context, msg *MessageRequest) error
	subscribe(ctx context.Context) (<-chan *MessageResponse, context.CancelFunc)
	CreateFromOpts(ctx context.Context) error
	Create(ctx context.Context, req *QueueCreateRequest) error
	Close() error
}

type apophis struct {
	ops    *options
	conn   *grpc.ClientConn
	client pb.PubSubServiceClient
}

func New(ops *options) ApophisInterface {
	if errs, ok := ops.Validate(); !ok {
		panic(errs)
	}

	var opts []grpc.DialOption
	if ops.insecured {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(nil)))
	}

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", ops.host, ops.port), opts...)
	if err != nil {
		panic(err)
	}

	return &apophis{
		conn:   conn,
		client: pb.NewPubSubServiceClient(conn),
		ops:    ops,
	}
}

func (a *apophis) Ping() (err error) {
	_, err = a.client.Ping(context.Background(), &pb.PingRequest{})
	return
}

func (a *apophis) Create(ctx context.Context, req *QueueCreateRequest) error {
	_, err := a.client.Create(ctx, req.PubRequest)
	return err
}

func (a *apophis) CreateFromOpts(ctx context.Context) error {
	_, err := a.client.Create(ctx, a.ops.GetPubRequest())
	return err
}

func (a *apophis) Publish(ctx context.Context, msg *MessageRequest) (err error) {

	if msg.Uniqid == "" {
		msg.Uniqid = a.ops.queueName
	} 

	_, err = a.client.Publish(ctx, msg.PubMessageRequest)
	if err != nil {
		fmt.Println("client.Publish err:")
	}
	return
}

func (a *apophis) subscribe(ctx context.Context) (<-chan *MessageResponse, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	response := make(chan *MessageResponse)
	go func() {
		defer func() {
			close(response)
		}()
		for {
			err := a.watching(ctx, response)
			if err != nil {
				fmt.Println("watching err:", err)
				fmt.Println("watching error status: ", status.Code(err))
				if status.Code(err) != codes.Unknown {
					time.Sleep(a.ops.reconnectInterval)
				} else {
					break
				}
			}
		}
	}()
	return response, cancel
}

func (a *apophis) read_messages(stream grpc.BidiStreamingClient[pb.SubscribeMessage, pb.SubscribeMessage], topic chan<- *MessageResponse, errCh chan<- error) {
	for {
		msg, err := stream.Recv()
		if err != nil {
			errCh <- err
			break
		}

		out := &MessageResponse{SubscribeMessage: msg}

		autoCommit := time.AfterFunc(a.ops.autoCommitTime, func() {
			msg.Commit = pb.MessageCommit_DISCARD
			stream.Send(msg)
		})

		out.OK = func() {
			autoCommit.Stop()
			msg.Commit = pb.MessageCommit_OK
			stream.Send(msg)
		}

		out.Retry = func() {
			autoCommit.Stop()
			msg.Commit = pb.MessageCommit_RETRY
			stream.Send(msg)
		}

		out.Ignore = func() {
			autoCommit.Stop()
			msg.Commit = pb.MessageCommit_DISCARD
			stream.Send(msg)
		}

		topic <- out
	}
}

func (a *apophis) watching(ctx context.Context, topic chan *MessageResponse) error {
	if err := a.Ping(); err != nil {
		return err
	}
	res, err := a.client.Subscribe(context.Background())
	if err != nil {
		return err
	}
	err = res.Send(&pb.SubscribeMessage{
		Sign: &pb.SubscribeRequest{
			Uniqid:      a.ops.queueName,
			Parallelism: int32(a.ops.consumerParralelism),
		},
	})

	if err != nil {
		res.CloseSend()
		return err
	}

	forward := make(chan *MessageResponse)
	errCh := make(chan error, 1)
	defer close(forward)
	go a.read_messages(res, forward, errCh)
	for {
		select {
		case <-ctx.Done():
			_ = res.Send(&pb.SubscribeMessage{
				UnSing: &pb.UnSubscribeRequest{
					Uniqid: a.ops.queueName,
				},
			})
			return fmt.Errorf("context canceled")
		case err := <-errCh:
			return err
		case msg := <-forward:
			topic <- msg
		}
	}
}

func (a *apophis) Close() error {
	return a.conn.Close()
}
