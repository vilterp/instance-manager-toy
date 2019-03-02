package manager

import "time"

type subID int

type mockOpLog struct {
	ops       []*Operation
	nextID    OpID
	opsByID   map[OpID]*Operation
	streams   map[subID]*mockStream
	nextSubID subID
}

var _ OpLog = &mockOpLog{}

func NewMockOpLog() *mockOpLog {
	return &mockOpLog{
		opsByID: map[OpID]*Operation{},
		ops:     nil,
		streams: map[subID]*mockStream{},
	}
}

func (ol *mockOpLog) publish(evt OpEvent) {
	for _, stream := range ol.streams {
		stream.c <- evt
	}
}

func (ol *mockOpLog) insertOp(opName string) *Operation {
	op := &Operation{
		ID:      ol.nextID,
		Started: time.Now(),
		Name:    opName,
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

func (ol *mockOpLog) Get(id OpID) *Operation {
	return ol.opsByID[id]
}

func (ol *mockOpLog) Tail() Stream {
	subID := ol.nextSubID
	stream := &mockStream{
		id:  subID,
		log: ol,
		c:   make(chan OpEvent),
	}
	ol.streams[subID] = stream
	return stream
}

func (ol *mockOpLog) unsubscribe(id subID) {
	delete(ol.streams, id)
}

func (ol *mockOpLog) Wait(id OpID) error {
	stream := ol.Tail()

	op := ol.opsByID[id]
	if op.Finished != nil {
		return op.Err
	}

	for evt := range stream.Events() {
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

type mockStream struct {
	log *mockOpLog
	id  subID
	c   chan OpEvent
}

var _ Stream = &mockStream{}

func (ms mockStream) Events() chan OpEvent {
	return ms.c
}

func (ms mockStream) Unsubscribe() {
	ms.log.unsubscribe(ms.id)
}
