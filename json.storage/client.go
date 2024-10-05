package jsonstorage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ninesbr/sheeps.toolkit.go/json.storage/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	tokenHeader = "x-access-key"
)

type JsonStorageInterface interface {
	Ping() error
	Close() error
    pushDocuments(data [][]byte) error
	getDocument(id string, value any) error
	getDocuments(query map[string]string, value any) error
	patchDocuments(data [][]byte) error
	deleteDocuments(id ...string) error
}

type storage struct {
	ops    *options
	conn   *grpc.ClientConn
	client pb.JsonStorageServiceClient
}

func New(ops *options) JsonStorageInterface {
	if err := ops.Validate(); err != nil {
		panic(err)
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

	return &storage{
		conn:   conn,
		client: pb.NewJsonStorageServiceClient(conn),
		ops:    ops,
	}
}

func (s *storage) metadata() metadata.MD {
	return metadata.New(map[string]string{
		tokenHeader: s.ops.accessKey,
	})
}

func (s *storage) Ping() error {
	_, err := s.client.Ping(context.Background(), &pb.PingRequest{})
	return err
}

func (s *storage) Close() error {
	return s.conn.Close()
}

func (s *storage) pushDocuments(data [][]byte) error {
	ctx := metadata.NewOutgoingContext(context.Background(), s.metadata())
	_, err := s.client.PushDocuments(ctx, &pb.PushDocsRequest{
		Documents: data,
	})
	return err
}

func (s *storage) getDocument(id string, value any) error {
	ctx := metadata.NewOutgoingContext(context.Background(), s.metadata())

	resp, err := s.client.GetDocument(ctx, &pb.GetDocRequest{
		UniqueId: id,
	})
	if err != nil {
		return err
	}
	return json.Unmarshal(resp.Document, value)
}

func (s *storage) getDocuments(query map[string]string, value any) error {
	ctx := metadata.NewOutgoingContext(context.Background(), s.metadata())

	resp, err := s.client.GetDocuments(ctx, &pb.GetDocsRequest{
		Query: query,
	})
	if err != nil {
		return err
	}

	return json.Unmarshal(resp.Documents, value)
}

func (s *storage) deleteDocuments(id ...string) error {
	ctx := metadata.NewOutgoingContext(context.Background(), s.metadata())

	_, err := s.client.DeleteDocuments(ctx, &pb.DeleteDocsRequest{
		UniqueIds: id,
	})
	return err
}

func (s *storage) patchDocuments(data [][]byte) error {
	ctx := metadata.NewOutgoingContext(context.Background(), s.metadata())
	_, err := s.client.PatchDocuments(ctx, &pb.PatchDocsRequest{
		Documents: data,
	})
	return err
}

