package apophis

import (
	"context"
)

type Subscribe struct {
	cli    ApophisInterface
	ctx    context.Context
	cancel context.CancelFunc
}

func NewSubscribe(ops *options) *Subscribe {
	sub := &Subscribe{
		cli: New(ops),
	}
	sub.ctx, sub.cancel = context.WithCancel(context.Background())
	return sub
}

func (s *Subscribe) Run(callback func(msg *MessageResponse)) {
	go func() {
		messages, cancel := s.cli.subscribe(s.ctx)
		defer cancel()
		for msg := range messages {
			callback(msg)
		}
	}()
}

func (s *Subscribe) Close() {
	s.cancel()
}
