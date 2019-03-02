package pure_manager

import (
	"fmt"
	"strings"
)

// action tree

type ActionNode interface {
	ActionNode()
	String() string
}

type DoNothing struct {
}

func (DoNothing) ActionNode()    {}
func (DoNothing) String() string { return "Nothing" }

type ActionLeaf struct {
	action Action
}

func (ActionLeaf) ActionNode() {}
func (al ActionLeaf) String() string {
	return fmt.Sprintf("%#v", al.action)
}

type Serial struct {
	Nodes []ActionNode
}

func (Serial) ActionNode() {}
func (s Serial) String() string {
	var actions []string
	for _, a := range s.Nodes {
		actions = append(actions, a.String())
	}
	return fmt.Sprintf("Ser(%s)", strings.Join(actions, ", "))
}

type Parallel struct {
	Nodes []ActionNode
}

func (Parallel) ActionNode() {}
func (p Parallel) String() string {
	var actions []string
	for _, a := range p.Nodes {
		actions = append(actions, a.String())
	}
	return fmt.Sprintf("Par(%s)", strings.Join(actions, ", "))
}

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
		action: a,
	}
}
