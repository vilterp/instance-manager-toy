package taskgraph

import (
	"github.com/cockroachlabs/instance_manager/pure_manager/actions"
	"github.com/google/uuid"
)

type Spec struct {
	Tasks map[TaskID]*TaskSpec
}

type TaskSpec struct {
	Action   actions.Action
	Upstream []TaskID
}

func NewSpec() *Spec {
	return &Spec{
		Tasks: map[TaskID]*TaskSpec{},
	}
}

func (s *Spec) ParIDs(list []TaskID) TaskID {
	id := TaskID(uuid.New())
	s.Tasks[id] = &TaskSpec{
		Action:   actions.DoNothing{},
		Upstream: list,
	}
	return id
}

func (s *Spec) Par(list []actions.Action) TaskID {
	var ids []TaskID
	for _, action := range list {
		id := TaskID(uuid.New())
		ts := &TaskSpec{
			Action: action,
		}
		s.Tasks[id] = ts
		ids = append(ids, id)
	}
	downstreamID := TaskID(uuid.New())
	downstream := &TaskSpec{
		Action:   actions.DoNothing{},
		Upstream: ids,
	}
	s.Tasks[downstreamID] = downstream
	return downstreamID
}

func (s *Spec) Ser(actions []actions.Action) TaskID {
	mostDownstream := TaskID(uuid.Nil)
	for _, action := range actions {
		id := TaskID(uuid.New())
		spec := &TaskSpec{
			Action: action,
		}
		if mostDownstream != TaskID(uuid.Nil) {
			spec.Upstream = []TaskID{mostDownstream}
		}
		mostDownstream = id
		s.Tasks[id] = spec
	}
	return mostDownstream
}
