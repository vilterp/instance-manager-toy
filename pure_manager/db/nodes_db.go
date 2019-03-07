package db

import "github.com/cockroachlabs/instance_manager/pure_manager/proto"

type SubID int

type NodeStateDB interface {
	List() []*proto.Node
	ListHealthy() []*proto.Node
	Stream() (SubID, chan *proto.NodeEvent)
	Unsubscribe(id SubID)

	Insert(*proto.Node)
	UpdateState(id proto.NodeID, newState proto.NodeState)
	//UpdateHealthStatus(id proto.NodeID, res HealthCheckResult)
}

type mockNodesDB struct {
	instancesByID map[proto.NodeID]*proto.Node
	instancesList []*proto.Node

	subs      map[SubID]chan *proto.NodeEvent
	nextSubID SubID
}

var _ NodeStateDB = &mockNodesDB{}

func NewEmptyMockInstancesDB() *mockNodesDB {
	return &mockNodesDB{
		instancesByID: map[proto.NodeID]*proto.Node{},
		subs:          map[SubID]chan *proto.NodeEvent{},
	}
}

func (m *mockNodesDB) Stream() (SubID, chan *proto.NodeEvent) {
	c := make(chan *proto.NodeEvent)
	subID := m.nextSubID
	m.subs[subID] = c
	m.nextSubID++
	return subID, c
}

func (m *mockNodesDB) Unsubscribe(id SubID) {
	delete(m.subs, id)
}

func (m *mockNodesDB) UpdateState(id proto.NodeID, newState proto.NodeState) {
	m.instancesByID[id].State = newState

	m.publish(&proto.NodeEvent{
		Event: &proto.NodeEvent_StateChanged_{
			StateChanged: &proto.NodeEvent_StateChanged{
				Id:       string(id),
				NewState: newState,
			},
		},
	})
}

func (m *mockNodesDB) Insert(n *proto.Node) {
	m.instancesByID[proto.NodeID(n.Id)] = n
	m.instancesList = append(m.instancesList, n)

	m.publish(&proto.NodeEvent{
		Event: &proto.NodeEvent_Started_{
			Started: &proto.NodeEvent_Started{
				Node: n,
			},
		},
	})
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

func (m *mockNodesDB) publish(evt *proto.NodeEvent) {
	for _, c := range m.subs {
		c <- evt
	}
}
