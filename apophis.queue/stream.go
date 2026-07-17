package apophis

import (
	"context"
	"sync"
	"time"

	"github.com/ninesbr/sheeps.toolkit.go/apophis.queue/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// serializedStreamSender aplica a regra de concorrência do gRPC: uma goroutine
// pode receber enquanto outra envia, mas dois Send simultâneos não são seguros.
type serializedStreamSender struct {
	mu     sync.Mutex
	stream grpc.BidiStreamingClient[pb.SubscribeMessage, pb.SubscribeMessage]
}

func (s *serializedStreamSender) Send(msg *pb.SubscribeMessage) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.stream.Send(msg)
}

// deliveryConfirmation concentra o estado de confirmação de uma mensagem.
// O sync.Once faz OK, Retry, Discard e o timeout competirem por uma única ação.
type deliveryConfirmation struct {
	ctx    context.Context
	sender *serializedStreamSender
	msg    *pb.SubscribeMessage
	errCh  chan<- error

	once sync.Once
	done chan struct{}
}

func newDeliveryConfirmation(
	ctx context.Context,
	sender *serializedStreamSender,
	msg *pb.SubscribeMessage,
	errCh chan<- error,
) *deliveryConfirmation {
	return &deliveryConfirmation{
		ctx:    ctx,
		sender: sender,
		// A confirmação mantém uma cópia privada. O callback pode consultar ou
		// alterar os headers expostos sem competir com a montagem do commit.
		msg:   proto.Clone(msg).(*pb.SubscribeMessage),
		errCh: errCh,
		done:  make(chan struct{}),
	}
}

func (c *deliveryConfirmation) confirm(commit pb.MessageCommit, headers map[string]string) {
	c.once.Do(func() {
		close(c.done)

		// Não alteramos a mensagem recebida: o gRPC pode manter referências à
		// mensagem enviada, e mutá-la depois de Send criaria outra corrida.
		ack := proto.Clone(c.msg).(*pb.SubscribeMessage)
		ack.Commit = commit
		if len(headers) > 0 {
			if ack.Headers == nil {
				ack.Headers = make(map[string]string, len(headers))
			}
			for key, value := range headers {
				ack.Headers[key] = value
			}
		}

		if err := c.sender.Send(ack); err != nil {
			// Os callbacks públicos não retornam erro. Encaminhamos a falha ao loop
			// interno para que ele encerre este stream e aplique a reconexão.
			select {
			case c.errCh <- err:
			case <-c.ctx.Done():
			default:
			}
		}
	})
}

func (c *deliveryConfirmation) cancel() {
	c.once.Do(func() {
		close(c.done)
	})
}

func (c *deliveryConfirmation) startTimeout(timeout time.Duration) {
	go func() {
		timer := time.NewTimer(timeout)
		defer timer.Stop()

		select {
		case <-timer.C:
			c.confirm(pb.MessageCommit_DISCARD, nil)
		case <-c.done:
		case <-c.ctx.Done():
			// O stream terminou; uma confirmação posterior não deve usar o stream
			// antigo nem competir com a próxima conexão.
			c.cancel()
		}
	}()
}
