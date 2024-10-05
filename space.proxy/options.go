package spaceproxy

import "errors"

type options struct {
	host             string
	port             int
	insecured        bool
	chunkSize        int
	uploadConcurrent int
}

func NewOptions(ops ...func(*options)) *options {
	svr := &options{
		chunkSize:        2 * 1024 * 1024, // 2MB
		uploadConcurrent: 1,
		insecured:        false,
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

func WithChunkSize(chunkSize int) func(*options) {
	return func(o *options) {
		o.chunkSize = chunkSize
	}
}

func WithUploadConcurrent(uploadConcurrent int) func(*options) {
	return func(o *options) {
		o.uploadConcurrent = uploadConcurrent
	}
}

func (o *options) Validate() (err error) {
	if o.host == "" {
		err = errors.Join(err, errors.New("host is empty"))
	}

	if o.port == 0 {
		err = errors.Join(err, errors.New("port is empty"))
	}

	if o.chunkSize == 0 {
		err = errors.Join(err, errors.New("chunkSize is empty"))
	}

	return
}
