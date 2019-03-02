package pure_manager

import (
	"fmt"
	"testing"
)

func TestDecide(t *testing.T) {
	st := StateDB{
		groupSpec: NewSpecDB(GroupSpec{}),
		instances: NewEmptyMockInstancesDB(),
		opLog:     NewMockOpLog(),
	}

	input := &CommandInput{
		Command: &UpdateSpec{
			NewSpec: GroupSpec{
				Version:      1,
				NumInstances: 3,
			},
		},
	}

	// TODO: make it easier to set up data driven tests

	actions := Decide(st, input)
	fmt.Printf("%#v\n", actions.String())

	t.Fatal()
}
