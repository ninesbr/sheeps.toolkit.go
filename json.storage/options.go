package jsonstorage

import (
	"errors"
)

type options struct {
	host      string
	port      int
	insecured bool
	accessKey string
}

func NewOptions(ops ...func(*options)) *options {
	svr := &options{
		insecured: false,
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

func WithAccessKey(accessKey string) func(*options) {
	return func(o *options) {
		o.accessKey = accessKey
	}
}

func (o *options) Validate() (err error) {
	if o.host == "" {
		err = errors.Join(err, errors.New("host is required"))
	}

	if o.port == 0 {
		err = errors.Join(err, errors.New("port is required"))
	}

	if o.accessKey == "" {
		err = errors.Join(err, errors.New("accessKey is required"))
	}

	return
}
