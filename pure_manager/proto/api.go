package proto

import (
	"fmt"
	"sort"
)

type Version int64

type NodeID string

type TaskID string

type idAndTS struct {
	spec *TaskSpec
	id   string
}

func (g *TaskGraphSpec) OrderedTaskSpecs() []idAndTS {
	var specs []idAndTS
	for tsid, ts := range g.Tasks {
		specs = append(specs, idAndTS{
			id:   tsid,
			spec: ts,
		})
	}
	sort.Slice(specs, func(i, j int) bool {
		return specs[i].id < specs[j].id
	})
	return specs
}

func (g *TaskGraphSpec) Print() {
	fmt.Println("ID\tprereqs\taction")
	for _, ts := range g.OrderedTaskSpecs() {
		fmt.Printf("%s\t%v\t%+v\n", ts.id, ts.spec.PrereqTaskIds, ts.spec.Action)
	}
}
