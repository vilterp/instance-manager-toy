package pure_manager

import "github.com/cockroachlabs/instance_manager/pure_manager/proto"

type Command interface {
	Command()
}

type UpdateSpec struct {
	NewSpec GroupSpec
}

func (us UpdateSpec) Command() {}

type KillInstance struct {
	ID proto.InstanceID
}

func (ki KillInstance) Command() {}
