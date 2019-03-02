package manager

import (
	"log"
	"testing"
)

func TestManager_Update(t *testing.T) {
	runnerLog := NewMockOpLog()
	runner := NewMockRunner(runnerLog, fastKnobs)

	tail := runnerLog.Tail()

	go func() {
		for opEvt := range tail.Events() {
			op := runnerLog.Get(opEvt.OpID())
			log.Printf("runner op: %T %d %#v", opEvt, op.ID, op.Name)
		}
	}()

	m := NewManager(GroupSpec{}, runner)
	err := m.Update(GroupSpec{
		Version:      1,
		NumInstances: 3,
	})
	if err != nil {
		t.Fatal(err)
	}

	m.WaitTilStable()

	log.Println("initial scale up succeeded; adding a node")

	err = m.Update(GroupSpec{
		Version:      1,
		NumInstances: 4,
	})
	if err != nil {
		t.Fatal(err)
	}

	m.WaitTilStable()

	expectedOps := []string{
		"start instance with spec {Version:1}",
		"start instance with spec {Version:1}",
		"start instance with spec {Version:1}",

		"start instance with spec {Version:1}",
	}

	compareOps(t, expectedOps, runnerLog.GetAll())
}
