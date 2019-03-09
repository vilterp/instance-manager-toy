package server

import (
	"github.com/cockroachlabs/instance_manager/proto"
	"github.com/cockroachlabs/instance_manager/taskgraph"
)

type builder struct {
	b *taskgraph.SpecBuilder
}

func newBuilder() *builder {
	return &builder{
		b: taskgraph.NewSpecBuilder(),
	}
}

// the dumbest possible strategy
func (b *builder) WipeAndRestart(nodes []*proto.Node, desiredSpec *proto.GroupSpec) *proto.TaskGraphSpec {
	w := b.Wipe(nodes)
	s := b.StartFromScratch(desiredSpec)
	b.b.SerIDs([]proto.TaskID{
		w,
		s,
	})
	return b.b.Build()
}

func (b *builder) Wipe(nodes []*proto.Node) proto.TaskID {
	var out []proto.TaskID
	for _, i := range nodes {
		out = append(out, b.b.Unit(&proto.Action{
			Action: &proto.Action_ShutDownNode{
				ShutDownNode: &proto.ShutDownNode{
					NodeId: i.Id,
				},
			},
		}))
	}
	return b.b.ParIDs("Wipe", out)
}

func (b *builder) StartFromScratch(spec *proto.GroupSpec) proto.TaskID {
	var out []proto.TaskID
	for i := int64(0); i < spec.NumInstances; i++ {
		out = append(out, b.b.Unit(&proto.Action{
			Action: &proto.Action_StartNode{
				StartNode: &proto.StartNode{
					Spec: &proto.NodeSpec{
						Version: spec.Version,
					},
				},
			},
		}))
	}
	return b.b.ParIDs("StartFromScratch", out)
}

//func RollingRestart(nodes []*Instance, newVersion Version) ActionNode {
//	var out []ActionNode
//	for _, i := range nodes {
//		if i.Version == newVersion {
//			continue
//		}
//		out = append(out, Unit(&RestartInstance{
//			ID:         i.ID,
//			NewVersion: newVersion,
//		}))
//	}
//	return Ser(out)
//}
//
//func ReplaceInstance(id InstanceID, spec GroupSpec) ActionNode {
//	return Par([]ActionNode{
//		Unit(ShutDownInstance{id}),
//		Unit(StartInstance{InstanceSpec{spec.Version}}),
//	})
//}
