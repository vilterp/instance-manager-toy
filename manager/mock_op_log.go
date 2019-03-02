package manager

import "time"

type mockOpLog struct {
	ops       []*Operation
	nextID    OpID
	opsByID   map[OpID]*Operation
	tailChans []chan OpEvent
}

var _ OpLog = &mockOpLog{}

func NewMockOpLog() *mockOpLog {
	return &mockOpLog{
		opsByID:   map[OpID]*Operation{},
		ops:       nil,
		tailChans: nil,
	}
}

func (ol *mockOpLog) publish(evt OpEvent) {
	for _, c := range ol.tailChans {
		c <- evt
	}
}

func (ol *mockOpLog) insertOp(opName string) *Operation {
	op := &Operation{
		ID:      ol.nextID,
		Started: time.Now(),
		Op:      opName,
	}
	ol.nextID++
	ol.ops = append(ol.ops, op)
	ol.opsByID[op.ID] = op
	return op
}

func (ol *mockOpLog) OpStarted(opName string) *Operation {
	op := ol.insertOp(opName)
	ol.publish(&OpStarted{
		Op: op,
	})
	return op
}

func (ol *mockOpLog) OpSucceeded(id OpID) {
	op := ol.opsByID[id]
	now := time.Now()
	op.Finished = &now
	ol.publish(&OpSucceeded{
		ID: id,
	})
}

func (ol *mockOpLog) OpFailed(id OpID, err error) {
	op := ol.opsByID[id]
	now := time.Now()
	op.Finished = &now
	op.Err = err
	ol.publish(&OpFailed{
		ID:  id,
		Err: err,
	})
}

func (ol *mockOpLog) GetAll() []*Operation {
	return ol.ops
}

func (ol *mockOpLog) Tail() chan OpEvent {
	c := make(chan OpEvent)
	// TODO: how does unsubscribing work?
	ol.tailChans = append(ol.tailChans, c)
	return c
}

func (ol *mockOpLog) Wait(id OpID) error {
	op := ol.opsByID[id]
	if op.Finished != nil {
		return op.Err
	}
	for evt := range ol.Tail() {
		switch t := evt.(type) {
		case *OpFailed:
			if t.ID == id {
				return t.Err
			}
		case *OpSucceeded:
			if t.ID == id {
				return nil
			}
		}
	}
	return nil
}
