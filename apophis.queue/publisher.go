package apophis

import "encoding/json"

type Publisher[T any] struct {
	cli ApophisInterface
}

func NewPublisher[T any](cli ApophisInterface) *Publisher[T] {
	return &Publisher[T]{
		cli: cli,
	}
}

func (p *Publisher[T]) Publish(msg *T, opsFunc ...func(*MessageRequestOptions)) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	req := &MessageRequest{
		ContentType: "application/json",
		Body:        data,
	}

	ops := &MessageRequestOptions{}
	for _, op := range opsFunc {
		op(ops)
	}

	if ops.ContentType != "" {
		req.ContentType = ops.ContentType
	}

    if ops.Headers != nil {
		req.Headers = ops.Headers
	}

    if ops.Tags != nil {
		req.Tags = ops.Tags
	}

    if ops.CustomID != "" {
		req.CustomID = ops.CustomID
	}

	if ops.TrackingID != "" {
		req.TrackingID = ops.TrackingID
	}

	return p.cli.publish(req)
}

