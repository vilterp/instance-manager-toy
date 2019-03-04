package taskgraph

import (
	"time"

	"github.com/cockroachlabs/instance_manager/pure_manager/actions"
	"github.com/google/uuid"
)

type TaskID uuid.UUID

func (t TaskID) String() string {
	return uuid.UUID(t).String()
}

// Supposed to be threadsafe
type StateDB interface {
	Insert(id TaskID, a actions.Action) TaskID
	AddDep(doFirst TaskID, thenDo TaskID)
	GetUnblockedTasks() []*Task
	List() []*Task

	MarkStarted(id TaskID)
	MarkSucceeded(id TaskID)
	MarkFailed(id TaskID, err error)

	// TODO: tail?
}

type TaskStatus string

const (
	StatusWaiting   TaskStatus = "WAITING"
	StatusRunning              = "RUNNING"
	StatusSucceeded            = "DONE"
	StatusFailed               = "FAILED"
)

type Task struct {
	ID         TaskID
	Action     actions.Action
	Status     TaskStatus
	StartedAt  time.Time
	FinishedAt time.Time
	// set if it failed
	Err error

	// uhh, this may need to be an interface eventually
}
