package server

import "github.com/cockroachlabs/instance_manager/pure_manager/proto"

type NodeStateDB interface {
	List() []*proto.Node
	ListHealthy() []*proto.Node

	Insert(*proto.Node)
	UpdateHealthStatus(id proto.NodeID, res HealthCheckResult)
}

// TODO: runner

type mockNodesDB struct {
	instancesByID map[proto.NodeID]*proto.Node
	instancesList []*proto.Node
}

var _ NodeStateDB = &mockNodesDB{}

func NewEmptyMockInstancesDB() *mockNodesDB {
	return &mockNodesDB{
		instancesByID: map[proto.NodeID]*proto.Node{},
	}
}

func (m *mockNodesDB) Insert(n *proto.Node) {
	m.instancesByID[proto.NodeID(n.Id)] = n
	m.instancesList = append(m.instancesList, n)
}

func (m *mockNodesDB) List() []*proto.Node {
	return m.instancesList
}

func (m *mockNodesDB) ListHealthy() []*proto.Node {
	var out []*proto.Node
	for _, i := range m.instancesList {
		out = append(out, i)
	}
	return out
}

func (m *mockNodesDB) UpdateHealthStatus(id proto.NodeID, res HealthCheckResult) {
	switch res.(type) {
	case HealthOk:
		m.instancesByID[id].State = proto.NodeState_NodeRunning
	case HealthErr:
		m.instancesByID[id].State = proto.NodeState_NodeUnhealthy
	}
}
