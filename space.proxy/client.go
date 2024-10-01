package spaceproxy

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/ninesbr/sheeps.toolkit.go/space.proxy/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type UploadRequest struct {
	Key         string
	ContentType string
	Extension   string
	Size        int64
}

type UploadResponse struct {
	*pb.PushRes
}

type SpaceInterface interface {
	Upload(ctx context.Context, info *UploadRequest, buf *bufio.Reader) (*UploadResponse, error)
	Close() error
}

type space struct {
	ops    *options
	conn   *grpc.ClientConn
	client pb.StorageCloudServiceClient
}

func New(ops *options) SpaceInterface {
	if errs, ok := ops.Validate(); !ok {
		panic(errs)
	}

	var opts []grpc.DialOption
	if ops.insecured {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(nil)))
	}

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", ops.host, ops.port), opts...)
	if err != nil {
		panic(err)
	}

	return &space{
		conn:   conn,
		client: pb.NewStorageCloudServiceClient(conn),
		ops:    ops,
	}
}

func (s *space) Upload(ctx context.Context, info *UploadRequest, buf *bufio.Reader) (*UploadResponse, error) {
	var (
		err    error
		res    *pb.PushRes
		stream grpc.ClientStreamingClient[pb.PushReq, pb.PushRes]
	)
	stream, err = s.client.Push(ctx)
	if err != nil {
		return nil, err
	}

	md := &pb.PushReq_Metadata{
		Metadata: &pb.Metadata{
			Key:         info.Key,
			ContentType: info.ContentType,
			Concurrent:  uint32(s.ops.uploadConcurrent),
			Extension:   info.Extension,
			Size:        uint64(info.Size),
		},
	}

	if err = stream.Send(&pb.PushReq{
		Data: md,
	}); err != nil {
		return nil, err
	}

	buffer := make([]byte, s.ops.chunkSize)

	for {
		n, err := buf.Read(buffer)
		if err != nil {
			if err == os.ErrClosed || err == io.EOF {
				break
			}
		}

		if n == 0 {
			break
		}

		if err = stream.Send(&pb.PushReq{
			Data: &pb.PushReq_Chunk{
				Chunk: buffer[:n],
			},
		}); err != nil {
			return nil, err
		}
	}

	res, err = stream.CloseAndRecv()
	if err != nil {
		return nil, err
	}

	return &UploadResponse{res}, nil
}

func (s *space) Close() error {
	return s.conn.Close()
}
