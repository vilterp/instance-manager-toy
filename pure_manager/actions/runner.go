package actions

import (
	"log"
	"sync"
	"time"
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

func (m *MockRunner) Run(a Action) error {
	log.Println("running", a.String())
	time.Sleep(1 * time.Second)
	m.mu.Lock()
	m.Log = append(m.Log, a)
	m.mu.Unlock()
	return nil
}
