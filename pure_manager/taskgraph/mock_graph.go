package taskgraph

import (
	"time"

	"github.com/cockroachlabs/instance_manager/pure_manager/actions"
	"github.com/google/uuid"
)

type MockGraphDB struct {
	tasks      map[TaskID]*Task
	downstream map[TaskID][]TaskID
	upstream   map[TaskID][]TaskID

	waitingTasks map[TaskID]struct{}
}

var _ StateDB = &MockGraphDB{}

func NewMockGraphDB(s *Spec) *MockGraphDB {
	g := &MockGraphDB{
		tasks:      map[TaskID]*Task{},
		downstream: map[TaskID][]TaskID{},
		upstream:   map[TaskID][]TaskID{},
	}
	for tID, taskSpec := range s.Tasks {
		g.Insert(taskSpec.Action)
		for _, upstream := range taskSpec.Upstream {
			g.AddDep(upstream, tID)
		}
	}
	return g
}

func (g *MockGraphDB) List() []*Task {
	var out []*Task
	for _, t := range g.tasks {
		out = append(out, t)
	}
	return out
}

func (g *MockGraphDB) Insert(a actions.Action) TaskID {
	id := TaskID(uuid.New())
	task := &Task{
		ID:     id,
		Action: a,
		Status: StatusWaiting,
	}
	g.tasks[id] = task
	g.upstream[id] = nil
	g.downstream[id] = nil
	g.waitingTasks[id] = struct{}{}
	return id
}

func (g *MockGraphDB) AddDep(doFirst TaskID, thenDo TaskID) {
	g.downstream[doFirst] = append(g.downstream[doFirst], thenDo)
	g.upstream[thenDo] = append(g.upstream[thenDo], doFirst)
}

func (g *MockGraphDB) GetUnblockedTasks() []*Task {
	var out []TaskID
	for id, task := range g.waitingTasks {
		upstreams := g.upstream[id]
		for _, upstreamID := range upstreams {
			if g.tasks[upstreamID].Status != StatusSucceeded {
				break
			}
			out = append(out, task)
		}
	}
	return out
}

func (g *MockGraphDB) MarkStarted(id TaskID) {
	g.tasks[id].Status = StatusRunning
	delete(g.waitingTasks, id)
}

func (g *MockGraphDB) MarkSucceeded(id TaskID) {
	t := g.tasks[id]
	t.FinishedAt = time.Now()
	t.Status = StatusSucceeded
}

func (g *MockGraphDB) MarkFailed(id TaskID, err error) {
	t := g.tasks[id]
	t.FinishedAt = time.Now()
	t.Status = StatusFailed
	t.Err = err
}

// TODO: getDot
