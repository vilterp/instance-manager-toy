package server

import (
	"fmt"

	"github.com/vilterp/instance-manager-toy/db"
	"github.com/vilterp/instance-manager-toy/proto"
)

func Decide(st *db.StateDB, input *proto.Input) (*proto.TaskGraphSpec, error) {
	switch tInput := input.Input.(type) {
	case *proto.Input_UpdateSpec:
		spec := tInput.UpdateSpec.Spec
		st.GroupSpec.UpdateSpec(spec)
		b := newBuilder()

		nodes := st.Nodes.List()

		specDelta := spec.NumInstances - int64(len(nodes))
		fmt.Println("delta", specDelta)
		if specDelta > 0 {
			b.StartNodes(specDelta, spec.Version)
			return b.Build(), nil
		} else if specDelta < 0 {
			_, err := b.KillSome(nodes, -specDelta)
			if err != nil {
				return nil, err
			}
			return b.Build(), nil
		} else {
			// TODO: update version and size at the same time
			b.RollingUpgrade(nodes, &proto.NodeSpec{
				Version: spec.Version,
			})
			return b.Build(), nil
		}
		return nil, fmt.Errorf("don't know how to upgrade")
	case *proto.Input_KillNode:
		panic("implement me")
	}
	panic(fmt.Sprintf("unknown input type: %T", input))
}
