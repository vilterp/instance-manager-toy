package manager

import "time"

type OpLog interface {
	OpLogPublisher
	OpLogReader
}

type OpLogPublisher interface {
	OpStarted(op string) *Operation
	OpSucceeded(id OpID)
	OpFailed(id OpID, err error)
}

type OpLogReader interface {
	GetAll() []*Operation
	Tail() chan OpEvent
	Wait(id OpID) error
}

type OpID int

type Operation struct {
	ID       OpID
	Op       string
	Started  time.Time
	Finished *time.Time
	Err      error // failed if this is not nil
}

func (op *Operation) IsDone() bool {
	return op.Finished != nil
}

type OpEvent interface {
	OpID() OpID
}

type OpStarted struct {
	Op *Operation
}

func (os *OpStarted) OpID() OpID {
	return os.Op.ID
}

type OpSucceeded struct {
	ID OpID
}

func (os *OpSucceeded) OpID() OpID {
	return os.ID
}

type OpFailed struct {
	ID  OpID
	Err error
}

func (of *OpFailed) OpID() OpID {
	return of.ID
}
