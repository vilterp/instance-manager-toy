package proto

import (
	"fmt"
	"sort"
	"strconv"
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
		it, err1 := strconv.Atoi(specs[i].id)
		jt, err2 := strconv.Atoi(specs[j].id)
		if err1 != nil {
			return false
		}
		if err2 != nil {
			return false
		}
		return it < jt
	})
	return specs
}

func (g *TaskGraphSpec) Print() {
	fmt.Println("ID\tprereqs\taction")
	for _, ts := range g.OrderedTaskSpecs() {
		fmt.Printf("%s\t%v\t%+v\n", ts.id, ts.spec.PrereqTaskIds, ts.spec.Action)
	}
}

func MkStartNode(spec *NodeSpec) *Action {
	return &Action{
		Action: &Action_StartNode{
			StartNode: &StartNode{
				Spec: spec,
			},
		},
	}
}

func MkShutDown(id NodeID) *Action {
	return &Action{
		Action: &Action_ShutDownNode{
			ShutDownNode: &ShutDownNode{
				NodeId: string(id),
			},
		},
	}
}
