package server

import (
	"fmt"

	"github.com/cockroachlabs/instance_manager/db"
	"github.com/cockroachlabs/instance_manager/proto"
)

func Decide(st *db.StateDB, input *proto.Input) (*proto.TaskGraphSpec, error) {
	switch tInput := input.Input.(type) {
	case *proto.Input_UpdateSpec:
		update := tInput.UpdateSpec
		st.GroupSpec.UpdateSpec(update.Spec)
		b := newBuilder()

		nodes := st.Nodes.List()

		specDelta := update.Spec.NumInstances - int64(len(nodes))
		fmt.Println("delta", specDelta)
		if specDelta > 0 {
			b.StartNodes(specDelta, update.Spec.Version)
			return b.Build(), nil
		} else if specDelta < 0 {
			_, err := b.KillSome(nodes, -specDelta)
			if err != nil {
				return nil, err
			}
			return b.Build(), nil
		}

		b.WipeAndRestart(st.Nodes.List(), st.GroupSpec.GetCurrent().Spec)
		return b.Build(), nil
	case *proto.Input_KillNode:
		panic("implement me")
	}
	panic(fmt.Sprintf("unknown input type: %T", input))
}
