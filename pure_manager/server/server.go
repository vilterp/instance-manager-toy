package server

import (
	"context"
	"log"

	"github.com/cockroachlabs/instance_manager/pure_manager/taskgraph"

	"github.com/cockroachlabs/instance_manager/pure_manager/actions"

	"github.com/cockroachlabs/instance_manager/pure_manager/proto"
)

type Server struct {
	db           *StateDB
	actionRunner actions.Runner
}

var _ proto.GroupManagerServer = &Server{}

func NewServer() *Server {
	return &Server{
		db:           NewStateDB(),
		actionRunner: actions.NewMockRunner(),
	}
}

func (s *Server) UpdateSpec(ctx context.Context, req *proto.UpdateSpecRequest) (*proto.UpdateSpecResponse, error) {
	graphSpec := Decide(s.db, &proto.Input{
		Input: &proto.Input_UpdateSpec{
			UpdateSpec: req,
		},
	})
	graph := s.db.taskGraphs.Insert(graphSpec)
	graphState := s.db.taskGraphs.GetState(TaskGraphID(graph.Id))
	runner := taskgraph.NewGraphRunner(graphState, s.actionRunner)
	go runner.Run()
	// TODO: run graph
	log.Printf("UpdateSpec: returning %#v", graph)
	return &proto.UpdateSpecResponse{
		GraphId: graph.Id,
	}, nil
}

func (s *Server) KillNode(context.Context, *proto.KillNodeRequest) (*proto.KillNodeResponse, error) {
	panic("implement me")
}

func (s *Server) GetCurrentSpec(context.Context, *proto.GetCurrentSpecRequest) (*proto.GetCurrentSpecResponse, error) {
	return &proto.GetCurrentSpecResponse{
		Spec: s.db.groupSpec.GetCurrent(),
	}, nil
}

func (s *Server) ListSpecs(context.Context, *proto.ListSpecsRequest) (*proto.ListSpecsResponse, error) {
	panic("implement me")
}

func (s *Server) ListNodes(context.Context, *proto.ListNodesRequest) (*proto.ListNodesResponse, error) {
	panic("implement me")
}

func (s *Server) StreamNodes(*proto.StreamNodesRequest, proto.GroupManager_StreamNodesServer) error {
	panic("implement me")
}

func (s *Server) ListTaskGraphs(context.Context, *proto.ListTaskGraphsRequest) (*proto.ListTaskGraphsResponse, error) {
	panic("implement me")
}

func (s *Server) StreamTaskGraphs(*proto.StreamTaskGraphsRequest, proto.GroupManager_StreamTaskGraphsServer) error {
	panic("implement me")
}

func (s *Server) GetTasks(context.Context, *proto.GetTasksRequest) (*proto.GetTasksResponse, error) {
	panic("implement me")
}

func (s *Server) StreamTasks(req *proto.StreamTasksRequest, srv proto.GroupManager_StreamTasksServer) error {
	st := s.db.taskGraphs.GetState(TaskGraphID(req.GraphId))
	if req.IncludeInitial {
		//if err := srv.Send(&proto.TaskEvent{
		//	Event: &proto.TaskEvent_Initial{
		//		Initial: &proto.TaskEvent_InitialState{
		//			Tasks: st.List(),
		//		},
		//	},
		//}); err != nil {
		//	return err
		//}
	}
	for evt := range st.Stream() {
		log.Println("evt:", evt, "srv:", srv)
		if err := srv.Send(evt); err != nil {
			return err
		}
	}
	return nil
}
