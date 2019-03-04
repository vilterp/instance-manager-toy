package pure_manager

import (
	"time"

	"github.com/cockroachlabs/instance_manager/pure_manager/taskgraph"
	"github.com/google/uuid"
)

type TaskGraphID uuid.UUID

// what task graphs are running
type TaskGraphsDB interface {
	Insert(g taskgraph.Spec)
	List() []*TaskGraph
	GetState(id TaskGraphID) taskgraph.StateDB
}

type TaskGraph struct {
	ID         uuid.UUID
	Spec       taskgraph.Spec
	StartedAt  time.Time
	FinishedAt time.Time
	Err        error
}

type MockTaskGraphsDB struct {
	graphs      map[TaskGraphID]*TaskGraph
	graphStates map[TaskGraphID]taskgraph.StateDB
}

func (g *MockTaskGraphsDB) List() []*TaskGraph {
	var out []*TaskGraph
	for _, tg := range g.graphs {
		out = append(out, tg)
	}
	return out
}

func (g *MockTaskGraphsDB) GetState(id TaskGraphID) taskgraph.StateDB {
	return g.graphStates[id]
}

var _ TaskGraphsDB = &MockTaskGraphsDB{}

func NewMockTaskGraphsDB() *MockTaskGraphsDB {
	return &MockTaskGraphsDB{
		graphs: map[TaskGraphID]*TaskGraph{},
	}
}

var _ TaskGraphsDB = &MockTaskGraphsDB{}

func (g *MockTaskGraphsDB) Insert(g taskgraph.Spec) {
	id := TaskGraphID(uuid.New())
	g.graphs[id] = taskgraph.NewMockGraphDB()
}
