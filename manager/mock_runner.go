package manager

import (
	"fmt"
	"time"
)

// could break out InstanceStore

type mockRunner struct {
	instancesByID map[InstanceID]*Instance
	instances     []*Instance
	nextID        InstanceID
	opLog         OpLog
}

var _ Runner = &mockRunner{}

func NewMockRunner(ol OpLog) *mockRunner {
	return &mockRunner{
		instancesByID: map[InstanceID]*Instance{},
		opLog:         ol,
	}
}

func (mr mockRunner) insertInstance(s InstanceSpec) *Instance {
	i := &Instance{
		ID:      mr.nextID,
		Version: s.Version,
		State:   StateStarting,
	}
	mr.nextID++
	mr.instances = append(mr.instances, i)
	mr.instancesByID[i.ID] = i
	return i
}

func (mr mockRunner) Start(spec InstanceSpec) (*Instance, *Operation, error) {
	i := mr.insertInstance(spec)
	op := mr.opLog.OpStarted(fmt.Sprintf("start instance with spec %v", spec))
	go func() {
		fmt.Println("starting up")
		time.Sleep(5 * time.Second)
		mr.opLog.OpSucceeded(op.ID)
		i.State = StateRunning
	}()
	return i, op, nil
}

func (mr mockRunner) ShutDown(id InstanceID) *Operation {
	i := mr.instancesByID[id]
	i.State = StateShuttingDown
	op := mr.opLog.OpStarted(fmt.Sprintf("shutting down instance %v", id))
	go func() {
		time.Sleep(5 * time.Second)
		i.State = StateRunning
		mr.opLog.OpSucceeded(op.ID)
	}()
	return op
}

func (mr mockRunner) GetInstance(id InstanceID) *Instance {
	return mr.instancesByID[id]
}

func (mr mockRunner) ListInstances() []*Instance {
	return mr.instances
}

func (mr mockRunner) GetOpLog() OpLogReader {
	return mr.opLog
}
