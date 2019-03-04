package pure_manager

import "github.com/cockroachlabs/instance_manager/pure_manager/base"

type InstanceStateDB interface {
	List() []*base.Instance
	ListHealthy() []*base.Instance

	Insert(*base.Instance)
	UpdateHealthStatus(id base.InstanceID, res HealthCheckResult)
}

// TODO: runner

type mockInstancesDB struct {
	instancesByID map[base.InstanceID]*base.Instance
	instancesList []*base.Instance
}

var _ InstanceStateDB = &mockInstancesDB{}

func NewEmptyMockInstancesDB() *mockInstancesDB {
	return &mockInstancesDB{
		instancesByID: map[base.InstanceID]*base.Instance{},
	}
}

func (m *mockInstancesDB) Insert(i *base.Instance) {
	m.instancesByID[i.ID] = i
	m.instancesList = append(m.instancesList, i)
}

func (m *mockInstancesDB) List() []*base.Instance {
	return m.instancesList
}

func (m *mockInstancesDB) ListHealthy() []*base.Instance {
	var out []*base.Instance
	for _, i := range m.instancesList {
		out = append(out, i)
	}
	return out
}

func (m *mockInstancesDB) UpdateHealthStatus(id base.InstanceID, res HealthCheckResult) {
	switch res.(type) {
	case HealthOk:
		m.instancesByID[id].State = base.StateRunning
	case HealthErr:
		m.instancesByID[id].State = base.StateUnhealthy
	}
}
