package manager

import (
	"fmt"
	"log"
)

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
		op := m.runner.ShutDown(r.ID)
		log.Printf("shut down instance %d: op %d", r.ID, op.ID)
	}

	for i := 0; i < newSpec.NumInstances; i++ {
		instance, op, err := m.runner.Start(InstanceSpec{
			Version: newSpec.Version,
		})
		if err != nil {
			// TODO: don't return an error, just keep retrying
			return nil
		}
		log.Printf("start up instance %d: op %d", instance.ID, op.ID)
	}

	return nil
}

// Diff returns nil if it's stable, or an error describing what's wrong.
func (m *Manager) Diff() error {
	return Matches(m.spec, m.runner.ListUpInstances())
}

func Matches(spec GroupSpec, instances []*Instance) error {
	if len(instances) != spec.NumInstances {
		return fmt.Errorf("want %d instances; have %v", spec.NumInstances, instances)
	}
	var wrongVersionIDs []InstanceID
	for _, instance := range instances {
		if instance.Version != spec.Version {
			wrongVersionIDs = append(wrongVersionIDs, instance.ID)
		}
	}
	if len(wrongVersionIDs) > 0 {
		return fmt.Errorf("instances %v are at the wrong version", wrongVersionIDs)
	}
	return nil
}

func (m *Manager) WaitTilStable() {
	diff := m.Diff()
	log.Printf("waiting til stable. diff: %v", diff)
	if diff == nil {
		return
	}
	stream := m.runner.GetOpLog().Tail()
	for op := range stream.Events() {
		diff := m.Diff()
		log.Printf("waiting til stable. op: %#v; diff: %v", op, diff)
		if diff == nil {
			stream.Unsubscribe()
			return
		}
	}
	return
}
