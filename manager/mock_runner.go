package manager

import (
	"fmt"
	"time"
)

type testingKnobs struct {
	instanceStartDur    DurDist
	instanceShutdownDur DurDist
}

var defaultKnobs = testingKnobs{
	instanceShutdownDur: DurDist{3 * time.Second, 3 * time.Second},
	instanceStartDur:    DurDist{3 * time.Second, 3 * time.Second},
}

var fastKnobs = testingKnobs{
	instanceShutdownDur: DurDist{100 * time.Millisecond, 100 * time.Millisecond},
	instanceStartDur:    DurDist{100 * time.Millisecond, 100 * time.Millisecond},
}

// TODO: break out InstanceStore
// in-memory vs DB

type mockRunner struct {
	instancesByID map[InstanceID]*Instance
	instances     []*Instance
	nextID        InstanceID
	opLog         OpLog
	testingKnobs  testingKnobs
}

var _ Runner = &mockRunner{}

func NewMockRunner(ol OpLog, knobs testingKnobs) *mockRunner {
	return &mockRunner{
		instancesByID: map[InstanceID]*Instance{},
		opLog:         ol,
		testingKnobs:  knobs,
	}
}

func (mr *mockRunner) insertInstance(s InstanceSpec) *Instance {
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

func (mr *mockRunner) Start(spec InstanceSpec) (*Instance, *Operation, error) {
	i := mr.insertInstance(spec)
	op := mr.opLog.OpStarted(fmt.Sprintf("start instance with spec %+v", spec))
	go func() {
		mr.testingKnobs.instanceStartDur.SleepRandom()
		mr.opLog.OpSucceeded(op.ID)
		i.State = StateRunning
	}()
	return i, op, nil
}

func (mr *mockRunner) ShutDown(id InstanceID) *Operation {
	i := mr.instancesByID[id]
	i.State = StateShuttingDown
	op := mr.opLog.OpStarted(fmt.Sprintf("shutting down instance %#v", id))
	go func() {
		mr.testingKnobs.instanceShutdownDur.SleepRandom()
		i.State = StateShutDown
		mr.opLog.OpSucceeded(op.ID)
	}()
	return op
}

func (mr *mockRunner) GetInstance(id InstanceID) *Instance {
	return mr.instancesByID[id]
}

func (mr *mockRunner) ListInstances() []*Instance {
	return mr.instances
}

func (mr *mockRunner) ListUpInstances() []*Instance {
	var out []*Instance
	for _, i := range mr.instances {
		if i.State == StateRunning {
			out = append(out, i)
		}
	}
	return out
}

func (mr *mockRunner) GetOpLog() OpLogReader {
	return mr.opLog
}
