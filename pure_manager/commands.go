package pure_manager

type Command interface {
	Command()
}

type UpdateSpec struct {
	NewSpec GroupSpec
}

func (us UpdateSpec) Command() {}

type KillInstance struct {
	ID InstanceID
}

func (ki KillInstance) Command() {}
