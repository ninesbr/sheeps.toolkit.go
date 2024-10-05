package apophis

import (
	"encoding/json"
)

type ConfirmDelivery struct {
	OK              func()
	Retry           func()
	RetryWithHeader func(header map[string]string)
	Discard         func()
}

type MessageResponse[T any] struct {
	ConfirmDelivery
	header  map[string]string
	body    []byte
	payload T
}

type MessageRequest struct {
	ContentType string
	Body        []byte
	Headers     map[string]string
	Tags        []string
	CustomID    string
	TrackingID  string
}

type MessageRequestOptions struct {
	ContentType string
	Headers     map[string]string
	Tags        []string
	CustomID    string
	TrackingID  string
}

func (res *MessageResponse[T]) GetHeader(key string) string {
	if res.header == nil {
		return ""
	}
	return res.header[key]
}

func (res *MessageResponse[T]) GetBody() []byte {
	return res.body
}

func (res *MessageResponse[T]) UnMarshalBody(v interface{}) error {
	if res.body == nil {
		return nil
	}
	err := json.Unmarshal(res.body, v)
	if err != nil {
		return err
	}
	return nil
}

func (res *MessageResponse[T]) GetHeaders() map[string]string {
	if res.header == nil {
		return map[string]string{}
	}
	return res.header
}

func (res *MessageResponse[T]) GetPayload() T {
	return res.payload
}

func WithHeaderRequest(key, value string) func(*MessageRequestOptions) {
	return func(o *MessageRequestOptions) {
		if o.Headers == nil {
			o.Headers = make(map[string]string)
		}
		o.Headers[key] = value
	}
}

func WithTagRequest(tag ...string) func(*MessageRequestOptions) {
	return func(o *MessageRequestOptions) {
		o.Tags = append(o.Tags, tag...)
	}
}

func WithCustomIDRequest(customID string) func(*MessageRequestOptions) {
	return func(o *MessageRequestOptions) {
		o.CustomID = customID
	}
}

func WithTrackingIDRequest(trackingID string) func(*MessageRequestOptions) {
	return func(o *MessageRequestOptions) {
		o.TrackingID = trackingID
	}
}

func WithContentTypeRequest(contentType string) func(*MessageRequestOptions) {
	return func(o *MessageRequestOptions) {
		o.ContentType = contentType
	}
}
