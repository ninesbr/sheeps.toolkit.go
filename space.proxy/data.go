package spaceproxy

import (
	"github.com/ninesbr/sheeps.toolkit.go/space.proxy/pb"
)

type CopyRequest struct {
	Key     string
	Uri     string
	Headers map[string]string
}

type CopyResponse struct {
	*pb.CopyFromRes
}

type UploadRequest struct {
	Key         string
	ContentType string
	Extension   string
	Size        int64
}

type UploadResponse struct {
	*pb.PushRes
}
