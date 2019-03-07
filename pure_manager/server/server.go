package server

import (
	"context"

	"github.com/cockroachlabs/instance_manager/pure_manager/actions"
	"github.com/cockroachlabs/instance_manager/pure_manager/db"
	"github.com/cockroachlabs/instance_manager/pure_manager/proto"
	"github.com/cockroachlabs/instance_manager/pure_manager/taskgraph"
)

type Server struct {
	db           *db.StateDB
	actionRunner actions.Runner
}

var _ proto.GroupManagerServer = &Server{}

func NewServer() *Server {
	db := db.NewStateDB()
	return &Server{
		db:           db,
		actionRunner: actions.NewMockRunner(db.Nodes),
	}
}

func (s *Server) UpdateSpec(ctx context.Context, req *proto.UpdateSpecRequest) (*proto.UpdateSpecResponse, error) {
	graphSpec := Decide(s.db, &proto.Input{
		Input: &proto.Input_UpdateSpec{
			UpdateSpec: req,
		},
	})
	graph := s.db.TaskGraphs.Insert(graphSpec)
	graphState := s.db.TaskGraphs.GetState(db.TaskGraphID(graph.Id))
	runner := taskgraph.NewGraphRunner(graphState, s.actionRunner)
	go func() {
		runner.Run()
		graphState.MarkGraphDone()
	}()
	return &proto.UpdateSpecResponse{
		Graph: graph,
	}, nil
}

func (s *Server) KillNode(context.Context, *proto.KillNodeRequest) (*proto.KillNodeResponse, error) {
	panic("implement me")
}

func (s *Server) GetCurrentSpec(context.Context, *proto.GetCurrentSpecRequest) (*proto.GetCurrentSpecResponse, error) {
	return &proto.GetCurrentSpecResponse{
		Spec: s.db.GroupSpec.GetCurrent(),
	}, nil
}

func (s *Server) ListSpecs(context.Context, *proto.ListSpecsRequest) (*proto.ListSpecsResponse, error) {
	panic("implement me")
}

func (s *Server) ListNodes(context.Context, *proto.ListNodesRequest) (*proto.ListNodesResponse, error) {
	panic("implement me")
}

func (s *Server) StreamNodes(req *proto.StreamNodesRequest, srv proto.GroupManager_StreamNodesServer) error {
	if req.IncludeInitial {
		if err := srv.Send(&proto.NodeEvent{
			Event: &proto.NodeEvent_InitialList_{
				InitialList: &proto.NodeEvent_InitialList{
					Nodes: s.db.Nodes.List(),
				},
			},
		}); err != nil {
			return err
		}
	}
	subID, c := s.db.Nodes.Stream()
	defer s.db.Nodes.Unsubscribe(subID)

	for evt := range c {
		if err := srv.Send(evt); err != nil {
			return err
		}
	}
	return nil
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
	st := s.db.TaskGraphs.GetState(db.TaskGraphID(req.GraphId))
	if req.IncludeInitial {
		for _, t := range st.List() {
			var b []byte
			if _, err := t.XXX_Marshal(b, true); err != nil {
				panic(err)
			}
		}
		if err := srv.Send(&proto.TaskEvent{
			Event: &proto.TaskEvent_Initial{
				Initial: &proto.TaskEvent_InitialState{
					Tasks: st.List(),
				},
			},
		}); err != nil {
			return err
		}
	}
	subID, c := st.Stream()
	defer st.Unsubscribe(subID)

	for evt := range c {
		if err := srv.Send(evt); err != nil {
			return err
		}
	}
	return nil
}
