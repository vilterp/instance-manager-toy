package taskgraph

import (
	"github.com/cockroachlabs/instance_manager/pure_manager/proto"
)

// Supposed to be threadsafe
type StateDB interface {
	Insert(id proto.TaskID, a *proto.Action) proto.TaskID
	AddDep(doFirst proto.TaskID, thenDo proto.TaskID)
	GetUnblockedTasks() []*proto.Task
	List() []*proto.Task
	Stream() chan *proto.TaskEvent

	MarkStarted(id proto.TaskID)
	MarkSucceeded(id proto.TaskID)
	MarkFailed(id proto.TaskID, err string)
	MarkGraphDone()

	// TODO: tail?
}
