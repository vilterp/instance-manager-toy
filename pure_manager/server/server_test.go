package server

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/cockroachlabs/instance_manager/pure_manager/proto"
	"google.golang.org/grpc/metadata"
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
	}, &mockTaskSrv{})
	if err3 != nil {
		log.Fatal(err3)
	}

	fmt.Println("updating again")
	_, err4 := s.UpdateSpec(ctx, &proto.UpdateSpecRequest{
		Spec: &proto.GroupSpec{
			NumInstances: 4,
		},
	})
	if err4 != nil {
		log.Fatal(err4)
	}

	err5 := s.StreamNodes(&proto.StreamNodesRequest{
		IncludeInitial: true,
	}, &mockNodeSrv{})
	if err5 != nil {
		log.Fatal(err5)
	}

	// TODO: get events
}

type mockTaskSrv struct {
}

func (s *mockTaskSrv) Send(evt *proto.TaskEvent) error {
	fmt.Println("sent", evt)
	return nil
}

func (mockTaskSrv) SetHeader(metadata.MD) error {
	panic("implement me")
}

func (mockTaskSrv) SendHeader(metadata.MD) error {
	panic("implement me")
}

func (mockTaskSrv) SetTrailer(metadata.MD) {
	panic("implement me")
}

func (mockTaskSrv) Context() context.Context {
	panic("implement me")
}

func (mockTaskSrv) SendMsg(m interface{}) error {
	panic("implement me")
}

func (mockTaskSrv) RecvMsg(m interface{}) error {
	panic("implement me")
}

var _ proto.GroupManager_StreamTasksServer = &mockTaskSrv{}

type mockNodeSrv struct{}

func (mockNodeSrv) Send(evt *proto.NodeEvent) error {
	fmt.Println("sent", evt)
	return nil
}

func (mockNodeSrv) SetHeader(metadata.MD) error {
	panic("implement me")
}

func (mockNodeSrv) SendHeader(metadata.MD) error {
	panic("implement me")
}

func (mockNodeSrv) SetTrailer(metadata.MD) {
	panic("implement me")
}

func (mockNodeSrv) Context() context.Context {
	panic("implement me")
}

func (mockNodeSrv) SendMsg(m interface{}) error {
	panic("implement me")
}

func (mockNodeSrv) RecvMsg(m interface{}) error {
	panic("implement me")
}

var _ proto.GroupManager_StreamNodesServer = &mockNodeSrv{}
