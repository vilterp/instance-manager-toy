package pure_manager

import "github.com/cockroachlabs/instance_manager/pure_manager/base"

type Command interface {
	Command()
}

type UpdateSpec struct {
	NewSpec GroupSpec
}

func (us UpdateSpec) Command() {}

type KillInstance struct {
	ID base.InstanceID
}

func (ki KillInstance) Command() {}
