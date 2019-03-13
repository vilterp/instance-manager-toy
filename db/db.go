package db

import "github.com/vilterp/instance-manager-toy/proto"

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
