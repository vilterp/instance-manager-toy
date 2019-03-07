package server

import "github.com/cockroachlabs/instance_manager/pure_manager/proto"

type StateDB struct {
	groupSpec  GroupSpecDB
	nodes      NodeStateDB
	taskGraphs TaskGraphsDB
}

func NewStateDB() *StateDB {
	return &StateDB{
		taskGraphs: NewMockTaskGraphsDB(),
		nodes:      NewEmptyMockInstancesDB(),
		groupSpec:  NewSpecDB(&proto.GroupSpec{}),
	}
}
