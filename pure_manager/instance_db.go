package pure_manager

import "github.com/cockroachlabs/instance_manager/pure_manager/proto"

type InstanceStateDB interface {
	List() []*proto.Instance
	ListHealthy() []*proto.Instance

	Insert(*proto.Instance)
	UpdateHealthStatus(id proto.InstanceID, res HealthCheckResult)
}

// TODO: runner

type mockInstancesDB struct {
	instancesByID map[proto.InstanceID]*proto.Instance
	instancesList []*proto.Instance
}

var _ InstanceStateDB = &mockInstancesDB{}

func NewEmptyMockInstancesDB() *mockInstancesDB {
	return &mockInstancesDB{
		instancesByID: map[proto.InstanceID]*proto.Instance{},
	}
}

func (m *mockInstancesDB) Insert(i *proto.Instance) {
	m.instancesByID[i.ID] = i
	m.instancesList = append(m.instancesList, i)
}

func (m *mockInstancesDB) List() []*proto.Instance {
	return m.instancesList
}

func (m *mockInstancesDB) ListHealthy() []*proto.Instance {
	var out []*proto.Instance
	for _, i := range m.instancesList {
		out = append(out, i)
	}
	return out
}

func (m *mockInstancesDB) UpdateHealthStatus(id proto.InstanceID, res HealthCheckResult) {
	switch res.(type) {
	case HealthOk:
		m.instancesByID[id].State = proto.StateRunning
	case HealthErr:
		m.instancesByID[id].State = proto.StateUnhealthy
	}
}
