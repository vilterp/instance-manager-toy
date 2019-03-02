package manager

import (
	"fmt"
	"testing"
)

func TestManager_Update(t *testing.T) {
	runnerLog := NewMockOpLog()
	runner := NewMockRunner(runnerLog)

	tail := runnerLog.Tail()

	go func() {
		for evt := range tail {
			fmt.Printf("runner op: %#v\n", evt)
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
}
