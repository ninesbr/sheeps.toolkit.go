package apophis

import (
	"context"
	"encoding/json"
)

type Subscribe[T any] struct {
	cli    ApophisInterface
	ctx    context.Context
	cancel context.CancelFunc
}

func NewSubscribe[T any](cli ApophisInterface) *Subscribe[T] {
	sub := &Subscribe[T]{
		cli: cli,
	}
	sub.ctx, sub.cancel = context.WithCancel(context.Background())
	return sub
}

func (s *Subscribe[T]) Run(callback func(msg *MessageResponse[T])) error {
	err := s.cli.Create()
	if err != nil {
		return err
	}
	go func() {
		messages, cancel := s.cli.subscribe(s.ctx)
		defer cancel()
		for msg := range messages {
			var ref T
			_ = json.Unmarshal(msg.GetBody(), &ref)
			msg.payload = ref
			callback(&MessageResponse[T]{
				ConfirmDelivery: msg.ConfirmDelivery,
				header:          msg.header,
				body:            msg.body,
				payload:         ref,
			})
		}
	}()
	return nil
}
