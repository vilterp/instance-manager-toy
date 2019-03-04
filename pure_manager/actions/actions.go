package actions

import (
	"fmt"

	"github.com/cockroachlabs/instance_manager/pure_manager/base"
)

// action tree

type ActionNode interface {
	ActionNode()
	String() string
}

// atomic actions

type Action interface {
	Action()
	String() string
}

type StartInstance struct {
	Spec base.InstanceSpec
}

func (StartInstance) Action()          {}
func (s StartInstance) String() string { return fmt.Sprintf("StartInstance(%v)", s.Spec) }

type ShutDownInstance struct {
	ID base.InstanceID
}

func (ShutDownInstance) Action()          {}
func (s ShutDownInstance) String() string { return fmt.Sprintf("ShutDownInstance(%v)", s.ID) }

type RestartInstance struct {
	ID         base.InstanceID
	NewVersion base.Version
}

func (RestartInstance) Action()          {}
func (r RestartInstance) String() string { return fmt.Sprintf("RestartInstance(%v)", r.ID) }

// TODO: region stuff
// TODO: restart at different size
