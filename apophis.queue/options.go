package apophis

import (
	"errors"
	"time"

	"github.com/ninesbr/sheeps.toolkit.go/apophis.queue/pb"
)

type options struct {
	host                string
	port                int
	insecured           bool
	reconnectInterval   time.Duration
	autoCommitTime      time.Duration
	consumerParralelism int

	queueName          string
	queueDurable       bool
	queueKeepMessages  bool
	queueTags          []string
	queueRetryInterval string
	queueRetryDuration string
}

func NewOptions(ops ...func(*options)) *options {
	svr := &options{
		reconnectInterval:   10 * time.Second,
		autoCommitTime:      5 * time.Second,
		consumerParralelism: 1,
	}
	for _, o := range ops {
		o(svr)
	}
	return svr
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

func WithInsecured(insecured bool) func(*options) {
	return func(o *options) {
		o.insecured = insecured
	}
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

func WithConsumerParralelism(parralelism int) func(*options) {
	return func(o *options) {
		o.consumerParralelism = parralelism
	}
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
	return
}
