package pure_manager

import (
	"fmt"

	"github.com/cockroachlabs/instance_manager/pure_manager/base"
)

type StateDB struct {
	groupSpec  GroupSpecDB
	instances  InstanceStateDB
	taskGraphs TaskGraphsDB
}

type Input interface {
	Input()
}

type CommandInput struct {
	Command Command
}

func (ci CommandInput) Input() {}

type HealthcheckInput struct {
	ID                base.InstanceID
	HealthcheckResult HealthCheckResult
}

func (hcr HealthcheckInput) Input() {}

type OpEventInput struct {
	OpEvent OpEvent
}

func (oi OpEventInput) Input() {}

func Decide(st StateDB, input Input) ActionNode {
	switch tInput := input.(type) {
	case *CommandInput:
		switch tCommand := tInput.Command.(type) {
		case *UpdateSpec:
			st.groupSpec.UpdateSpec(tCommand.NewSpec)
			return WipeAndRestart(st.instances.List(), st.groupSpec.GetCurrent())
		case *KillInstance:
			return ReplaceInstance(tCommand.ID, st.groupSpec.GetCurrent())
		}
	case *HealthcheckInput:
		st.instances.UpdateHealthStatus(tInput.ID, tInput.HealthcheckResult)
		return DoNothing{}
	case *OpEventInput:
		// TODO: replace with list of outstanding actions
		//   + action tree interpreter
		switch tEvt := tInput.OpEvent.(type) {
		case *OpStarted:
			st.opLog.OpStarted(tEvt.Name)
			return DoNothing{}
		case *OpSucceeded:
			st.opLog.OpSucceeded(tEvt.ID)
			// TODO: ok great, but what do we do about it?
			// have to update outstanding action trees
			return DoNothing{}
		case *OpFailed:
			st.opLog.OpFailed(tEvt.ID, tEvt.Err)
			return DoNothing{}
		}
	}
	panic(fmt.Sprintf("unknown input type: %T", input))
}
