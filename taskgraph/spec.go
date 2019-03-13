package taskgraph

import (
	"fmt"

	"github.com/vilterp/instance-manager-toy/proto"
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
		Action:        nothing(desc),
		PrereqTaskIds: ids,
	}
	return id
}

type TaskChain struct {
	Head proto.TaskID
	Tail proto.TaskID
}

func (b *SpecBuilder) SerIDs(desc string, list []proto.TaskID) TaskChain {
	var start proto.TaskID
	var mostDownstream proto.TaskID
	for i, t := range list {
		if i != 0 {
			b.AddDep(mostDownstream, t)
		} else {
			start = t
		}
		mostDownstream = t
	}
	return TaskChain{
		Head: start,
		Tail: mostDownstream,
	}
}

func (b *SpecBuilder) SerChains(chains []TaskChain) TaskChain {
	start := chains[0].Head
	mostDownstream := chains[0].Tail
	for _, c := range chains {
		if c.Tail == mostDownstream {
			continue
		}
		b.AddDep(mostDownstream, c.Head)
		mostDownstream = c.Tail
	}
	return TaskChain{
		Head: start,
		Tail: mostDownstream,
	}
}

func (b *SpecBuilder) AddDep(doFirst proto.TaskID, thenDo proto.TaskID) {
	downstream := b.spec.Tasks[string(thenDo)]
	downstream.PrereqTaskIds = append(downstream.PrereqTaskIds, string(doFirst))
}

func nothing(desc string) *proto.Action {
	return &proto.Action{
		Action: &proto.Action_DoNothing{
			DoNothing: &proto.DoNothing{
				Description: desc,
			},
		},
	}
}
