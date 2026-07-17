package apophis

import (
	"errors"
	"time"

	"github.com/ninesbr/sheeps.toolkit.go/apophis.queue/pb"
)

const maxInt32 = int64(1<<31 - 1)

type options struct {
	host                string
	port                int
	insecure            bool
	requestTimeout      time.Duration
	publishTimeout      time.Duration
	reconnectInterval   time.Duration
	autoCommitTime      time.Duration
	consumerParallelism int

	queueName          string
	queueDurable       bool
	queueKeepMessages  bool
	queueTags          []string
	queueRetryInterval string
	queueRetryDuration string
}

func NewOptions(ops ...func(*options)) *options {
	svr := &options{
		requestTimeout:      10 * time.Second,
		publishTimeout:      10 * time.Second,
		reconnectInterval:   10 * time.Second,
		autoCommitTime:      5 * time.Second,
		consumerParallelism: 1,
	}
	for _, o := range ops {
		o(svr)
	}
	return svr
}

// WithRequestTimeout limita as operações de controle Ping, Create e Drop.
// Publish possui um timeout próprio configurado por WithPublishTimeout.
func WithRequestTimeout(timeout time.Duration) func(*options) {
	return func(o *options) {
		o.requestTimeout = timeout
	}
}

// WithPublishTimeout limita uma chamada unary de Publish. O contexto não faz
// parte de ApophisInterface, portanto o prazo é aplicado dentro do cliente.
func WithPublishTimeout(timeout time.Duration) func(*options) {
	return func(o *options) {
		o.publishTimeout = timeout
	}
}

func WithHost(host string) func(*options) {
	return func(o *options) {
		o.host = host
	}
}

func WithPort(port int) func(*options) {
	return func(o *options) {
		o.port = port
	}
}

// WithInsecure define se a conexão deve usar credenciais de transporte sem
// TLS. O valor padrão é false, portanto conexões seguras são usadas por padrão.
func WithInsecure(insecure bool) func(*options) {
	return func(o *options) {
		o.insecure = insecure
	}
}

// WithInsecured mantém compatibilidade com a grafia antiga.
// Deprecated: use WithInsecure.
func WithInsecured(insecure bool) func(*options) {
	return WithInsecure(insecure)
}

func WithReconnectInterval(interval time.Duration) func(*options) {
	return func(o *options) {
		o.reconnectInterval = interval
	}
}

func WithAutoCommitTime(time time.Duration) func(*options) {
	return func(o *options) {
		o.autoCommitTime = time
	}
}

// WithConsumerParallelism informa ao servidor quantos consumidores da fila
// devem alimentar esta assinatura. O valor não controla a concorrência dos
// callbacks no cliente.
func WithConsumerParallelism(parallelism int) func(*options) {
	return func(o *options) {
		o.consumerParallelism = parallelism
	}
}

// WithConsumerParralelism mantém compatibilidade com a grafia antiga.
// Deprecated: use WithConsumerParallelism.
func WithConsumerParralelism(parallelism int) func(*options) {
	return WithConsumerParallelism(parallelism)
}

func WithQueueName(queueName string) func(*options) {
	return func(o *options) {
		o.queueName = queueName
	}
}

func WithQueueDurable(queueDurable bool) func(*options) {
	return func(o *options) {
		o.queueDurable = queueDurable
	}
}

func WithQueueKeepMessages(queueKeepMessages bool) func(*options) {
	return func(o *options) {
		o.queueKeepMessages = queueKeepMessages
	}
}

func WithQueueTags(queueTags ...string) func(*options) {
	return func(o *options) {
		o.queueTags = queueTags
	}
}

func WithQueueRetryInterval(queueRetryInterval string) func(*options) {
	return func(o *options) {
		o.queueRetryInterval = queueRetryInterval
	}
}

func WithQueueRetryDuration(queueRetryDuration string) func(*options) {
	return func(o *options) {
		o.queueRetryDuration = queueRetryDuration
	}
}

func (o *options) GetPubRequest() *pb.PubRequest {
	return &pb.PubRequest{
		Uniqid:        o.queueName,
		Durable:       o.queueDurable,
		KeepMessages:  o.queueKeepMessages,
		Tags:          o.queueTags,
		RetryInterval: o.queueRetryInterval,
		RetryDuration: o.queueRetryDuration,
	}
}

func (o *options) Validate() (err error) {
	if o.host == "" {
		err = errors.Join(err, errors.New("host is empty"))
	}
	if o.port == 0 {
		err = errors.Join(err, errors.New("port is empty"))
	}
	if o.queueName == "" {
		err = errors.Join(err, errors.New("queue name is empty"))
	}
	if o.requestTimeout <= 0 {
		err = errors.Join(err, errors.New("request timeout must be greater than zero"))
	}
	if o.publishTimeout <= 0 {
		err = errors.Join(err, errors.New("publish timeout must be greater than zero"))
	}
	if o.reconnectInterval <= 0 {
		err = errors.Join(err, errors.New("reconnect interval must be greater than zero"))
	}
	if o.autoCommitTime <= 0 {
		err = errors.Join(err, errors.New("auto commit time must be greater than zero"))
	}
	if o.consumerParallelism <= 0 {
		err = errors.Join(err, errors.New("consumer parallelism must be greater than zero"))
	} else if int64(o.consumerParallelism) > maxInt32 {
		err = errors.Join(err, errors.New("consumer parallelism exceeds int32"))
	}
	return
}
