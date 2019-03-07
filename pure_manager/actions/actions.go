package actions

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/cockroachlabs/instance_manager/pure_manager/db"
	"github.com/cockroachlabs/instance_manager/pure_manager/proto"
	"github.com/cockroachlabs/instance_manager/pure_manager/util"
)

type Runner interface {
	Run(a *proto.Action) error
}

type MockRunner struct {
	mu    sync.Mutex
	Log   []*proto.Action
	nodes db.NodeStateDB
}

func NewMockRunner(nodes db.NodeStateDB) *MockRunner {
	return &MockRunner{
		nodes: nodes,
	}
}

var _ Runner = &MockRunner{}

var sleepTimeDist = util.DurDist{Base: 1 * time.Second, Variance: 1 * time.Second}

func (m *MockRunner) Run(a *proto.Action) error {
	log.Println("running", a.String())
	m.mu.Lock()
	m.Log = append(m.Log, a)
	m.mu.Unlock()

	switch tAction := a.Action.(type) {
	case *proto.Action_StartNode:
		log.Println("starting node:", tAction)
		id := uuid.New().String()
		m.nodes.Insert(&proto.Node{
			Id:      id,
			Version: tAction.StartNode.Spec.Version,
			State:   proto.NodeState_NodeStarting,
		})
		sleepTimeDist.SleepRandom()
		m.nodes.UpdateState(proto.NodeID(id), proto.NodeState_NodeRunning)
	case *proto.Action_ShutDownNode:
		log.Println("shutting down node:", tAction)
		id := tAction.ShutDownNode.NodeId
		m.nodes.UpdateState(proto.NodeID(id), proto.NodeState_NodeShuttingDown)
		sleepTimeDist.SleepRandom()
		m.nodes.UpdateState(proto.NodeID(id), proto.NodeState_NodeShutDown)
	}

	return nil
}
