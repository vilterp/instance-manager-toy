package actions

import (
	"log"
	"sync"
	"time"

	"github.com/cockroachlabs/instance_manager/pure_manager/util"
)

type Runner interface {
	Run(a Action) error
}

type MockRunner struct {
	mu  sync.Mutex
	Log []Action
}

func NewMockRunner() *MockRunner {
	return &MockRunner{}
}

var _ Runner = &MockRunner{}

var dist = util.DurDist{Base: 1 * time.Second, Variance: 1 * time.Second}

func (m *MockRunner) Run(a Action) error {
	log.Println("running", a.String())
	m.mu.Lock()
	m.Log = append(m.Log, a)
	m.mu.Unlock()
	dist.SleepRandom()
	return nil
}
