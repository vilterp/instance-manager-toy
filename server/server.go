package server

import (
	"context"

	"github.com/vilterp/instance-manager-toy/actions"
	"github.com/vilterp/instance-manager-toy/db"
	"github.com/vilterp/instance-manager-toy/proto"
	"github.com/vilterp/instance-manager-toy/taskgraph"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *Server) UpdateSpec(ctx context.Context, req *proto.UpdateSpecReq) (*proto.UpdateSpecResp, error) {
	graphSpec, err := Decide(s.db, &proto.Input{
		Input: &proto.Input_UpdateSpec{
			UpdateSpec: req,
		},
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid spec:", err)
	}
	graphSpec.Print()
	graph := s.db.TaskGraphs.Insert(graphSpec)
	graphState := s.db.TaskGraphs.GetState(db.TaskGraphID(graph.Id))
	runner := taskgraph.NewRunner(graphState, s.actionRunner)
	go func() {
		runner.Run()
		s.db.TaskGraphs.MarkDone(db.TaskGraphID(graph.Id))
		graphState.MarkGraphDone()
	}()
	return &proto.UpdateSpecResp{
		Graph: graph,
	}, nil
}

func (s *Server) KillNode(context.Context, *proto.KillNodeReq) (*proto.KillNodeResp, error) {
	panic("implement me")
}

func (s *Server) GetCurrentSpec(context.Context, *proto.GetCurrentSpecReq) (*proto.GetCurrentSpecResp, error) {
	return &proto.GetCurrentSpecResp{
		Spec: s.db.GroupSpec.GetCurrent(),
	}, nil
}

func (s *Server) ListSpecs(context.Context, *proto.ListSpecsReq) (*proto.ListSpecsResp, error) {
	panic("implement me")
}

func (s *Server) ListNodes(context.Context, *proto.ListNodesReq) (*proto.ListNodesResp, error) {
	return &proto.ListNodesResp{
		Instances: s.db.Nodes.List(),
	}, nil
}

func (s *Server) StreamNodes(req *proto.StreamNodesReq, srv proto.GroupManager_StreamNodesServer) error {
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

func (s *Server) ListTaskGraphs(context.Context, *proto.ListTaskGraphsReq) (*proto.ListTaskGraphsResp, error) {
	return &proto.ListTaskGraphsResp{
		TaskGraphs: s.db.TaskGraphs.List(),
	}, nil
}

func (s *Server) StreamTaskGraphs(*proto.StreamTaskGraphsReq, proto.GroupManager_StreamTaskGraphsServer) error {
	panic("implement me")
}

func (s *Server) GetTaskGraph(ctx context.Context, req *proto.GetTaskGraphReq) (*proto.GetTaskGraphResp, error) {
	g, ok := s.db.TaskGraphs.Get(db.TaskGraphID(req.Id))
	if !ok {
		return nil, status.Error(codes.NotFound, "no graph with that id")
	}
	return &proto.GetTaskGraphResp{
		Graph: g,
	}, nil
}

func (s *Server) GetTasks(context.Context, *proto.GetTasksReq) (*proto.GetTasksResp, error) {
	panic("implement me")
}

func (s *Server) StreamTasks(req *proto.StreamTasksReq, srv proto.GroupManager_StreamTasksServer) error {
	st := s.db.TaskGraphs.GetState(db.TaskGraphID(req.GraphId))
	if st == nil {
		return status.Error(codes.NotFound, "no graph with that id")
	}
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
