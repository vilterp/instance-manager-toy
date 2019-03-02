package manager

import (
	"fmt"
	"testing"
)

func TestManager_Update(t *testing.T) {
	runnerLog := NewMockOpLog()
	runner := NewMockRunner(runnerLog)

	go func() {
		for evt := range runnerLog.Tail() {
			fmt.Println("runner op:", evt)
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
