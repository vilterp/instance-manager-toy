package taskgraph

import (
	"fmt"
	"log"

	"github.com/cockroachlabs/instance_manager/pure_manager/proto"
	"github.com/golang/protobuf/ptypes"
)

type MockGraphDB struct {
	tasks      map[proto.TaskID]*proto.Task
	downstream map[proto.TaskID][]proto.TaskID
	upstream   map[proto.TaskID][]proto.TaskID

	waitingTasks map[proto.TaskID]struct{}

	eventSubs map[int]chan *proto.TaskEvent
	nextSubID int
}

func (g *MockGraphDB) Print() {
	for _, t := range g.tasks {
		fmt.Println(t)
		up := g.upstream[proto.TaskID(t.Id)]
		down := g.downstream[proto.TaskID(t.Id)]
		fmt.Println("\tup:", len(up), up)
		fmt.Println("\tdown:", len(down), down)
	}
}

var _ StateDB = &MockGraphDB{}

func NewMockGraphDB(s *proto.TaskGraphSpec) *MockGraphDB {
	g := &MockGraphDB{
		tasks:        map[proto.TaskID]*proto.Task{},
		downstream:   map[proto.TaskID][]proto.TaskID{},
		upstream:     map[proto.TaskID][]proto.TaskID{},
		waitingTasks: map[proto.TaskID]struct{}{},
		eventSubs:    map[int]chan *proto.TaskEvent{},
	}
	for tID, taskSpec := range s.Tasks {
		g.Insert(proto.TaskID(tID), taskSpec.Action)
		for _, upstream := range taskSpec.PrereqTaskIds {
			g.AddDep(proto.TaskID(upstream), proto.TaskID(tID))
		}
	}
	return g
}

func (g *MockGraphDB) List() []*proto.Task {
	var out []*proto.Task
	for _, t := range g.tasks {
		out = append(out, t)
	}
	return out
}

func (g *MockGraphDB) Insert(id proto.TaskID, a *proto.Action) proto.TaskID {
	task := &proto.Task{
		Id:     string(id),
		Action: a,
		State:  proto.TaskState_TaskWaiting,
	}
	g.tasks[id] = task
	g.upstream[id] = nil
	g.downstream[id] = nil
	g.waitingTasks[id] = struct{}{}
	return id
}

func (g *MockGraphDB) AddDep(doFirst proto.TaskID, thenDo proto.TaskID) {
	g.downstream[doFirst] = append(g.downstream[doFirst], thenDo)
	g.upstream[thenDo] = append(g.upstream[thenDo], doFirst)
}

func (g *MockGraphDB) GetUnblockedTasks() []*proto.Task {
	var out []*proto.Task
	for id := range g.waitingTasks {
		upstreams := g.upstream[id]
		blocked := false
		for _, upstreamID := range upstreams {
			if g.tasks[upstreamID].State != proto.TaskState_TaskSucceeded {
				blocked = true
				break
			}
		}
		if !blocked {
			out = append(out, g.tasks[id])
		}
	}
	return out
}

func (g *MockGraphDB) MarkStarted(id proto.TaskID) {
	t := g.tasks[id]
	t.State = proto.TaskState_TaskRunning
	t.StartedAt = ptypes.TimestampNow()
	delete(g.waitingTasks, id)

	g.publish(&proto.TaskEvent{
		Event: &proto.TaskEvent_Started_{
			Started: &proto.TaskEvent_Started{
				Id: string(id),
			},
		},
	})
}

func (g *MockGraphDB) MarkSucceeded(id proto.TaskID) {
	t := g.tasks[id]
	t.FinishedAt = ptypes.TimestampNow()
	t.State = proto.TaskState_TaskSucceeded

	g.publish(&proto.TaskEvent{
		Event: &proto.TaskEvent_Succeeded_{
			Succeeded: &proto.TaskEvent_Succeeded{
				Id: string(id),
			},
		},
	})
}

func (g *MockGraphDB) MarkFailed(id proto.TaskID, err string) {
	t := g.tasks[id]
	t.FinishedAt = ptypes.TimestampNow()
	t.State = proto.TaskState_TaskFailed
	t.Error = err

	g.publish(&proto.TaskEvent{
		Event: &proto.TaskEvent_Failed_{
			Failed: &proto.TaskEvent_Failed{
				Id:    string(id),
				Error: err,
			},
		},
	})
}

func (g *MockGraphDB) publish(evt *proto.TaskEvent) {
	log.Println("publishing", evt)
	for _, c := range g.eventSubs {
		// TODO: one of these could block the rest...
		c <- evt
	}
	// TODO: close the conns when the last task succeeds
}

func (g *MockGraphDB) Stream() chan *proto.TaskEvent {
	c := make(chan *proto.TaskEvent)
	g.eventSubs[g.nextSubID] = c
	g.nextSubID++
	return c
}

// TODO: getDot
