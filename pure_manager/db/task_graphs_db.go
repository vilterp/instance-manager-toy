package db

import (
	"github.com/cockroachlabs/instance_manager/pure_manager/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
)

type TaskGraphID string

// what task graphs are running
type TaskGraphsDB interface {
	Insert(g *proto.TaskGraphSpec) *proto.TaskGraph
	List() []*proto.TaskGraph
	GetState(id TaskGraphID) TasksDB
}

type MockTaskGraphsDB struct {
	graphs      map[TaskGraphID]*proto.TaskGraph
	graphStates map[TaskGraphID]TasksDB
}

func (g *MockTaskGraphsDB) List() []*proto.TaskGraph {
	var out []*proto.TaskGraph
	for _, tg := range g.graphs {
		out = append(out, tg)
	}
	return out
}

func (g *MockTaskGraphsDB) GetState(id TaskGraphID) TasksDB {
	return g.graphStates[id]
}

var _ TaskGraphsDB = &MockTaskGraphsDB{}

func NewMockTaskGraphsDB() *MockTaskGraphsDB {
	return &MockTaskGraphsDB{
		graphs:      map[TaskGraphID]*proto.TaskGraph{},
		graphStates: map[TaskGraphID]TasksDB{},
	}
}

var _ TaskGraphsDB = &MockTaskGraphsDB{}

func (g *MockTaskGraphsDB) Insert(spec *proto.TaskGraphSpec) *proto.TaskGraph {
	id := TaskGraphID(uuid.New().String())
	graph := &proto.TaskGraph{
		Id:        string(id),
		Spec:      spec,
		State:     proto.TaskGraphState_TaskGraphWaiting,
		CreatedAt: ptypes.TimestampNow(),
	}
	g.graphs[id] = graph
	g.graphStates[id] = NewMockTasksDB(spec)
	return graph
}
