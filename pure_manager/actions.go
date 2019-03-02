package pure_manager

// action tree

type ActionNode interface {
	ActionNode()
}

type DoNothing struct {
}

func (DoNothing) ActionNode() {}

type ActionLeaf struct {
	A Action
}

func (ActionLeaf) ActionNode() {}

type Serial struct {
	Nodes []ActionNode
}

func (Serial) ActionNode() {}

type Parallel struct {
	Nodes []ActionNode
}

func (Parallel) ActionNode() {}

// atomic actions

type Action interface {
	Action()
}

type StartInstance struct {
	Spec InstanceSpec
}

func (StartInstance) Action() {}

type ShutDownInstance struct {
	ID InstanceID
}

func (ShutDownInstance) Action() {}

type RestartInstance struct {
	ID         InstanceID
	NewVersion Version
}

func (RestartInstance) Action() {}

// TODO: region stuff
// TODO: restart at different size

// Helpers

func Par(an []ActionNode) ActionNode {
	return &Parallel{
		Nodes: an,
	}
}

func Ser(an []ActionNode) ActionNode {
	return &Serial{
		Nodes: an,
	}
}

func Unit(a Action) ActionNode {
	return ActionLeaf{
		A: a,
	}
}
