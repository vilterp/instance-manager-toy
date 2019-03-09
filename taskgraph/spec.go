package taskgraph

import (
	"fmt"

	"github.com/cockroachlabs/instance_manager/proto"
)

type SpecBuilder struct {
	spec   *proto.TaskGraphSpec
	nextID int
}

func NewSpecBuilder() *SpecBuilder {
	return &SpecBuilder{
		spec: &proto.TaskGraphSpec{
			Tasks: map[string]*proto.TaskSpec{},
		},
	}
}

func (b *SpecBuilder) getNextID() proto.TaskID {
	id := b.nextID
	b.nextID++
	return proto.TaskID(fmt.Sprintf("%d", id))
}

func (b *SpecBuilder) Build() *proto.TaskGraphSpec {
	return b.spec
}

func (b *SpecBuilder) Unit(a *proto.Action) proto.TaskID {
	id := b.getNextID()
	b.spec.Tasks[string(id)] = &proto.TaskSpec{
		Action: a,
	}
	return id
}

func (b *SpecBuilder) ParIDs(desc string, list []proto.TaskID) proto.TaskID {
	id := b.getNextID()
	var ids []string
	for _, tid := range list {
		ids = append(ids, string(tid))
	}
	b.spec.Tasks[string(id)] = &proto.TaskSpec{
		Action: &proto.Action{
			Action: &proto.Action_DoNothing{
				DoNothing: &proto.DoNothing{
					Description: desc,
				},
			},
		},
		PrereqTaskIds: ids,
	}
	return id
}

func (b *SpecBuilder) SerIDs(list []proto.TaskID) proto.TaskID {
	mostDownstream := proto.TaskID("")
	for _, t := range list {
		if mostDownstream != proto.TaskID("") {
			b.AddDep(mostDownstream, t)
		}
		mostDownstream = t
	}
	return mostDownstream
}

func (b *SpecBuilder) AddDep(doFirst proto.TaskID, thenDo proto.TaskID) {
	downstream := b.spec.Tasks[string(thenDo)]
	downstream.PrereqTaskIds = append(downstream.PrereqTaskIds, string(doFirst))
}
