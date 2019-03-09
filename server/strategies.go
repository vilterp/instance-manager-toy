package server

import (
	"fmt"

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

func (b *builder) Build() *proto.TaskGraphSpec {
	return b.b.Build()
}

// the dumbest possible strategy
func (b *builder) WipeAndRestart(nodes []*proto.Node, desiredSpec *proto.GroupSpec) proto.TaskID {
	w := b.Wipe(nodes)
	s := b.StartFromScratch(desiredSpec)
	return b.b.SerIDs([]proto.TaskID{
		w,
		s,
	})
}

func (b *builder) KillSome(nodes []*proto.Node, n int64) (proto.TaskID, error) {
	if int64(len(nodes)) < n {
		return "", fmt.Errorf("not enough nodes to shut down")
	}
	var out []proto.TaskID
	for i := int64(0); i < n; i++ {
		node := nodes[i]
		out = append(out, b.b.Unit(&proto.Action{
			Action: &proto.Action_ShutDownNode{
				ShutDownNode: &proto.ShutDownNode{
					NodeId: node.Id,
				},
			},
		}))
	}
	return b.b.ParIDs("KillSome", out), nil
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
	return b.StartNodes(spec.NumInstances, spec.Version)
}

func (b *builder) StartNodes(n int64, v int64) proto.TaskID {
	var out []proto.TaskID
	for i := int64(0); i < n; i++ {
		out = append(out, b.b.Unit(&proto.Action{
			Action: &proto.Action_StartNode{
				StartNode: &proto.StartNode{
					Spec: &proto.NodeSpec{
						Version: int64(v),
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
