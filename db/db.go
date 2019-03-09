package db

import "github.com/cockroachlabs/instance_manager/proto"

type StateDB struct {
	GroupSpec  GroupSpecDB
	Nodes      NodeStateDB
	TaskGraphs TaskGraphsDB
}

func NewStateDB() *StateDB {
	return &StateDB{
		TaskGraphs: NewMockTaskGraphsDB(),
		Nodes:      NewEmptyMockInstancesDB(),
		GroupSpec:  NewSpecDB(&proto.GroupSpec{}),
	}
}
