package apophis

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/ninesbr/sheeps.toolkit.go/apophis.queue/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type ApophisInterface interface {
	Ping() error
	publish(msg *MessageRequest) error
	subscribe(ctx context.Context) (<-chan *MessageResponse[any], context.CancelFunc)
	Drop(keepMessagesRead bool) error
	Create() error
	Close() error
}

type apophis struct {
	ops    *options
	conn   *grpc.ClientConn
	client pb.PubSubServiceClient
}

func New(ops *options) ApophisInterface {
	if err := ops.Validate(); err != nil {
		panic(err)
	}

	var opts []grpc.DialOption
	if ops.insecure {
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
	ctx, cancel := a.requestContext(context.Background())
	defer cancel()

	_, err = a.client.Ping(ctx, &pb.PingRequest{})
	return
}

func (a *apophis) Create() error {
	ctx, cancel := a.requestContext(context.Background())
	defer cancel()

	_, err := a.client.Create(ctx, a.ops.GetPubRequest())
	return err
}

func (a *apophis) Drop(keepMessagesRead bool) error {
	ctx, cancel := a.requestContext(context.Background())
	defer cancel()

	_, err := a.client.Drop(ctx, &pb.DropRequest{
		Uniqid:           a.ops.queueName,
		KeepMessagesRead: keepMessagesRead,
	})
	return err
}

func (a *apophis) requestContext(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, a.ops.requestTimeout)
}

func (a *apophis) publish(msg *MessageRequest) (err error) {
	if msg == nil {
		return errors.New("message request is nil")
	}

	// ApophisInterface não recebe context no Publish. Aplicamos o timeout aqui
	// para que uma chamada unary nunca aguarde indefinidamente pelo servidor.
	ctx, cancel := context.WithTimeout(context.Background(), a.ops.publishTimeout)
	defer cancel()

	req := &pb.PubMessageRequest{
		ContentType: msg.ContentType,
		Uniqid:      a.ops.queueName,
		Headers:     msg.Headers,
		Body:        msg.Body,
		Tags:        msg.Tags,
		CustomID:    msg.CustomID,
		TrackingID:  msg.TrackingID,
	}
	if msg.ForceCreate {
		// ForceCreate é propositalmente opt-in para não mudar a política de
		// criação de filas dos publishers existentes.
		req.ForceCreate = a.ops.GetPubRequest()
	}

	_, err = a.client.Publish(ctx, req)
	return
}

func (a *apophis) subscribe(ctx context.Context) (<-chan *MessageResponse[any], context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	response := make(chan *MessageResponse[any])
	go func() {
		defer close(response)

		for {
			err := a.watching(ctx, response)

			// O cancelamento solicitado pelo consumidor sempre vence qualquer
			// política de reconexão.
			if ctx.Err() != nil {
				return
			}
			if !shouldReconnect(err) {
				if err != nil {
					fmt.Println("watching stopped:", err)
				}
				return
			}

			fmt.Println("watching reconnecting after error:", err)
			if !waitForReconnect(ctx, a.ops.reconnectInterval) {
				return
			}
		}
	}()
	return response, cancel
}

func shouldReconnect(err error) bool {
	if err == nil {
		return false
	}
	// Erros puros de context não são necessariamente convertidos em status
	// pelo gRPC e poderiam aparecer como Unknown.
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	// Recv retorna EOF quando o servidor encerra o stream com status OK. Para
	// um consumidor contínuo, isso significa abrir uma nova assinatura.
	if errors.Is(err, io.EOF) {
		return true
	}

	switch status.Code(err) {
	case codes.Unavailable,
		codes.ResourceExhausted,
		codes.Aborted,
		codes.Internal,
		codes.Unknown,
		codes.DeadlineExceeded:
		return true
	default:
		// Canceled, InvalidArgument, NotFound, PermissionDenied,
		// Unauthenticated e demais erros permanentes não entram em loop.
		return false
	}
}

func waitForReconnect(ctx context.Context, interval time.Duration) bool {
	timer := time.NewTimer(interval)
	defer timer.Stop()

	select {
	case <-timer.C:
		return true
	case <-ctx.Done():
		return false
	}
}

func (a *apophis) readMessages(ctx context.Context, stream grpc.BidiStreamingClient[pb.SubscribeMessage, pb.SubscribeMessage], sender *serializedStreamSender, topic chan<- *MessageResponse[any], errCh chan<- error) {
	for {
		msg, err := stream.Recv()
		if err != nil {
			select {
			case errCh <- err:
			case <-ctx.Done():
			}
			return
		}

		out := &MessageResponse[any]{
			header: msg.Headers,
			body:   msg.Body,
		}

		confirmation := newDeliveryConfirmation(ctx, sender, msg, errCh)

		out.OK = func() {
			confirmation.confirm(pb.MessageCommit_OK, nil)
		}

		out.Retry = func() {
			confirmation.confirm(pb.MessageCommit_RETRY, nil)
		}

		out.RetryWithHeader = func(header map[string]string) {
			confirmation.confirm(pb.MessageCommit_RETRY, header)
		}

		out.Discard = func() {
			confirmation.confirm(pb.MessageCommit_DISCARD, nil)
		}

		// A entrega também precisa respeitar o cancelamento. Sem este select, um
		// consumidor lento poderia manter a goroutine presa mesmo após Stop/Close.
		select {
		case topic <- out:
			// O prazo começa após o handoff; tempo bloqueado aguardando o
			// consumidor não deve contar como tempo de processamento.
			confirmation.startTimeout(a.ops.autoCommitTime)
		case <-ctx.Done():
			confirmation.cancel()
			return
		}
	}
}

func (a *apophis) watching(ctx context.Context, topic chan *MessageResponse[any]) error {
	streamCtx, cancelStream := context.WithCancel(ctx)
	defer cancelStream()

	// O mesmo contexto controla desde o teste de conexão até o Recv do stream.
	// Assim, cancelar a assinatura libera também os recursos internos do gRPC.
	pingCtx, cancelPing := a.requestContext(streamCtx)
	_, err := a.client.Ping(pingCtx, &pb.PingRequest{})
	cancelPing()
	if err != nil {
		return err
	}
	res, err := a.client.Subscribe(streamCtx)
	if err != nil {
		return err
	}

	errCh := make(chan error, 1)
	sender := &serializedStreamSender{stream: res}
	err = sender.Send(&pb.SubscribeMessage{
		Sign: &pb.SubscribeRequest{
			Uniqid:      a.ops.queueName,
			Parallelism: int32(a.ops.consumerParallelism),
		},
	})

	if err != nil {
		res.CloseSend()
		return err
	}

	go a.readMessages(streamCtx, res, sender, topic, errCh)

	select {
	case <-ctx.Done():
		// Não enviamos UnSing aqui: o contexto já encerra o stream inteiro e é a
		// forma recomendada pelo gRPC de liberar Recv e recursos associados.
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func (a *apophis) Close() error {
	return a.conn.Close()
}
