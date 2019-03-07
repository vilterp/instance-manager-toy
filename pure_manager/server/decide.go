package server

import (
	"fmt"

	"github.com/cockroachlabs/instance_manager/pure_manager/proto"
)

func Decide(st *StateDB, input *proto.Input) *proto.TaskGraphSpec {
	switch tInput := input.Input.(type) {
	case *proto.Input_UpdateSpec:
		update := tInput.UpdateSpec
		st.groupSpec.UpdateSpec(update.Spec)
		b := newBuilder()
		return b.WipeAndRestart(st.nodes.List(), st.groupSpec.GetCurrent().Spec)
	case *proto.Input_KillNode:
		panic("implement me")
	}
	panic(fmt.Sprintf("unknown input type: %T", input))
}
