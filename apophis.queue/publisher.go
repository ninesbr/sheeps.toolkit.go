package apophis

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"strings"
)

type Publisher[T any] struct {
	cli ApophisInterface
}

func NewPublisher[T any](cli ApophisInterface) *Publisher[T] {
	return &Publisher[T]{
		cli: cli,
	}
}

func (p *Publisher[T]) Publish(msg *T, opsFunc ...func(*MessageRequestOptions)) error {
	if p == nil || p.cli == nil {
		return errors.New("publisher client is nil")
	}
	if msg == nil {
		// json.Marshal aceitaria este valor e publicaria o JSON `null`, o que
		// normalmente representa uma mensagem inválida para o consumidor.
		return errors.New("message is nil")
	}

	ops := &MessageRequestOptions{}
	for index, op := range opsFunc {
		if op == nil {
			return fmt.Errorf("message request option %d is nil", index)
		}
		op(ops)
	}

	contentType := "application/json"
	if ops.ContentType != "" {
		if err := validateJSONContentType(ops.ContentType); err != nil {
			return err
		}
		contentType = ops.ContentType
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal publish message: %w", err)
	}

	req := &MessageRequest{
		ContentType: contentType,
		Body:        data,
		Headers:     ops.Headers,
		Tags:        ops.Tags,
		CustomID:    ops.CustomID,
		TrackingID:  ops.TrackingID,
		ForceCreate: ops.ForceCreate,
	}

	if err := p.cli.publish(req); err != nil {
		return fmt.Errorf("publish message: %w", err)
	}
	return nil
}

func validateJSONContentType(contentType string) error {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return fmt.Errorf("invalid content type %q: %w", contentType, err)
	}

	mediaType = strings.ToLower(mediaType)
	if mediaType != "application/json" && !strings.HasSuffix(mediaType, "+json") {
		// Este publisher não recebe bytes prontos: seu corpo sempre vem de
		// json.Marshal e não pode ser anunciado como protobuf, texto etc.
		return fmt.Errorf("content type %q is not compatible with JSON", contentType)
	}
	return nil
}
