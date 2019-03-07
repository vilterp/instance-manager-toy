package actions

import (
	"fmt"

	"github.com/cockroachlabs/instance_manager/pure_manager/proto"
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

type DoNothing struct {
}

func (DoNothing) Action()        {}
func (DoNothing) String() string { return "DoNothing" }

type StartInstance struct {
	Spec proto.InstanceSpec
}

func (StartInstance) Action()          {}
func (s StartInstance) String() string { return fmt.Sprintf("StartInstance(%v)", s.Spec) }

type ShutDownInstance struct {
	ID proto.InstanceID
}

func (ShutDownInstance) Action()          {}
func (s ShutDownInstance) String() string { return fmt.Sprintf("ShutDownInstance(%v)", s.ID) }

type RestartInstance struct {
	ID         proto.InstanceID
	NewVersion proto.Version
}

func (RestartInstance) Action()          {}
func (r RestartInstance) String() string { return fmt.Sprintf("RestartInstance(%v)", r.ID) }

// TODO: region stuff
// TODO: restart at different size
