package pure_manager

import (
	"time"

	"github.com/google/uuid"
)

type OpLog interface {
	OpLogPublisher
	OpLogReader
}

type OpLogPublisher interface {
	OpStarted(op string) *Operation
	OpSucceeded(id OpID)
	OpFailed(id OpID, err error)
}

type Stream interface {
	Events() chan OpEvent
	Unsubscribe()
}

type OpLogReader interface {
	List() []*Operation
	Get(id OpID) *Operation
	Tail() Stream
	Wait(id OpID) error
}

type OpID uuid.UUID

type Operation struct {
	ID       OpID
	Name     string
	Started  time.Time
	Finished *time.Time
	Err      error // failed if this is not nil
}

type OpEvent interface {
	OpID() OpID
}

type OpStarted struct {
	ID   OpID
	Name string
}

func (os *OpStarted) OpID() OpID {
	return os.ID
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

type mockOpLog struct {
	opsByID map[OpID]*Operation
	opsList []*Operation
}

var _ OpLog = &mockOpLog{}

func NewMockOpLog() *mockOpLog {
	return &mockOpLog{
		opsByID: map[OpID]*Operation{},
	}
}

func (m *mockOpLog) OpStarted(name string) *Operation {
	op := &Operation{
		ID:      OpID(uuid.New()),
		Started: time.Now(),
		Name:    name,
	}
	m.opsByID[op.ID] = op
	m.opsList = append(m.opsList, op)
	return op
}

func (m *mockOpLog) OpSucceeded(id OpID) {
	now := time.Now()
	m.opsByID[id].Finished = &now
}

func (m *mockOpLog) OpFailed(id OpID, err error) {
	now := time.Now()
	op := m.opsByID[id]
	op.Finished = &now
	op.Err = err
}

func (m *mockOpLog) List() []*Operation {
	return m.opsList
}

func (m *mockOpLog) Get(id OpID) *Operation {
	return m.opsByID[id]
}

func (m *mockOpLog) Tail() Stream {
	panic("implement me")
}

func (m *mockOpLog) Wait(id OpID) error {
	panic("implement me")
}
