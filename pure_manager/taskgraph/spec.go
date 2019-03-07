package taskgraph

import (
	"github.com/cockroachlabs/instance_manager/pure_manager/proto"
	"github.com/google/uuid"
)

type SpecBuilder struct {
	spec *proto.TaskGraphSpec
}

func NewSpecBuilder() *SpecBuilder {
	return &SpecBuilder{
		spec: &proto.TaskGraphSpec{
			Tasks: map[string]*proto.TaskSpec{},
		},
	}
}

func (b *SpecBuilder) Build() *proto.TaskGraphSpec {
	return b.spec
}

func (b *SpecBuilder) Unit(a *proto.Action) proto.TaskID {
	id := proto.TaskID(uuid.New().String())
	b.spec.Tasks[string(id)] = &proto.TaskSpec{
		Action: a,
	}
	return id
}

func (b *SpecBuilder) ParIDs(list []proto.TaskID) proto.TaskID {
	id := proto.TaskID(uuid.New().String())
	var ids []string
	for _, tid := range list {
		ids = append(ids, string(tid))
	}
	b.spec.Tasks[string(id)] = &proto.TaskSpec{
		Action: &proto.Action{
			Action: &proto.Action_DoNothing{},
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
