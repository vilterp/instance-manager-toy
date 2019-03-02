package manager

import "fmt"

type Manager struct {
	spec        GroupSpec
	specHistory []GroupSpec
	runner      Runner
	// TODO: something about current operation
}

type InstanceSpec struct {
	Version int
}

type GroupSpec struct {
	NumInstances int
	Version      int
}

func NewManager(spec GroupSpec, runner Runner) *Manager {
	return &Manager{
		spec:   spec,
		runner: runner,
	}
}

func (m *Manager) Update(newSpec GroupSpec) error {
	m.specHistory = append(m.specHistory, m.spec)
	m.spec = newSpec
	// TODO: run stuff

	for _, r := range m.runner.ListInstances() {
		fmt.Println("shut down")
		m.runner.ShutDown(r.ID)
	}

	for i := 0; i < newSpec.NumInstances; i++ {
		fmt.Println("start up")
		if _, _, err := m.runner.Start(InstanceSpec{
			Version: newSpec.Version,
		}); err != nil {
			// TODO: don't return an error, just keep retrying
			return nil
		}
	}

	return nil
}

// Diff returns nil if it's stable, or an error describing what's wrong.
func (m *Manager) Diff() error {
	return Matches(m.spec, m.runner.ListInstances())
}

func Matches(spec GroupSpec, instances []*Instance) error {
	if len(instances) != spec.NumInstances {
		return fmt.Errorf("want %d instancesByID; have %d", spec.NumInstances, len(instances))
	}
	var wrongVersionIDs []InstanceID
	for _, instance := range instances {
		if instance.Version != spec.Version {
			wrongVersionIDs = append(wrongVersionIDs, instance.ID)
		}
	}
	if len(wrongVersionIDs) > 0 {
		return fmt.Errorf("instancesByID %v are at the wrong version", wrongVersionIDs)
	}
	return nil
}

func (m *Manager) WaitTilStable() {
	for range m.runner.GetOpLog().Tail() {
		if m.Diff() == nil {
			return
		}
	}
	return
}
