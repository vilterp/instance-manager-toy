package pure_manager

import "github.com/google/uuid"

type InstanceState string

const (
	StateRunning      InstanceState = "RUNNING"
	StateStarting                   = "STARTING"
	StateUnhealthy                  = "UNHEALTHY"
	StateShuttingDown               = "SHUTTING_DOWN"
	StateShutDown                   = "SHUT_DOWN"
)

type InstanceStateDB interface {
	List() []*Instance
	ListHealthy() []*Instance

	Insert(*Instance)
	UpdateHealthStatus(id InstanceID, res HealthCheckResult)
}

type InstanceID uuid.UUID

type Version int

type Instance struct {
	ID      InstanceID
	Version Version
	State   InstanceState
}

type InstanceSpec struct {
	Version Version
}

// TODO: runner

type mockInstancesDB struct {
	instancesByID map[InstanceID]*Instance
	instancesList []*Instance
}

var _ InstanceStateDB = &mockInstancesDB{}

func NewEmptyMockInstancesDB() *mockInstancesDB {
	return &mockInstancesDB{
		instancesByID: map[InstanceID]*Instance{},
	}
}

func (m *mockInstancesDB) Insert(i *Instance) {
	m.instancesByID[i.ID] = i
	m.instancesList = append(m.instancesList, i)
}

func (m *mockInstancesDB) List() []*Instance {
	return m.instancesList
}

func (m *mockInstancesDB) ListHealthy() []*Instance {
	var out []*Instance
	for _, i := range m.instancesList {
		out = append(out, i)
	}
	return out
}

func (m *mockInstancesDB) UpdateHealthStatus(id InstanceID, res HealthCheckResult) {
	switch res.(type) {
	case HealthOk:
		m.instancesByID[id].State = StateRunning
	case HealthErr:
		m.instancesByID[id].State = StateUnhealthy
	}
}
