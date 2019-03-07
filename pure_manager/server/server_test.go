package server

import (
	"context"
	"fmt"
	"log"
	"testing"

	"google.golang.org/grpc/metadata"

	"github.com/cockroachlabs/instance_manager/pure_manager/proto"
)

func TestServer(t *testing.T) {
	s := NewServer()

	ctx := context.Background()
	resp, err := s.GetCurrentSpec(ctx, &proto.GetCurrentSpecRequest{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("get current:", resp)

	fmt.Println("update:")
	resp2, err2 := s.UpdateSpec(ctx, &proto.UpdateSpecRequest{
		Spec: &proto.GroupSpec{
			NumInstances: 3,
		},
	})
	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Println(resp2.Graph)

	fmt.Println("stream tasks:")
	err3 := s.StreamTasks(&proto.StreamTasksRequest{
		GraphId:        resp2.Graph.Id,
		IncludeInitial: true,
	}, &mockSrv{})
	if err3 != nil {
		log.Fatal(err3)
	}

	// TODO: get events
}

type mockSrv struct {
}

func (s *mockSrv) Send(evt *proto.TaskEvent) error {
	fmt.Println("sent", evt)
	return nil
}

func (mockSrv) SetHeader(metadata.MD) error {
	panic("implement me")
}

func (mockSrv) SendHeader(metadata.MD) error {
	panic("implement me")
}

func (mockSrv) SetTrailer(metadata.MD) {
	panic("implement me")
}

func (mockSrv) Context() context.Context {
	panic("implement me")
}

func (mockSrv) SendMsg(m interface{}) error {
	panic("implement me")
}

func (mockSrv) RecvMsg(m interface{}) error {
	panic("implement me")
}

var _ proto.GroupManager_StreamTasksServer = &mockSrv{}
