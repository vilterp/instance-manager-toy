package base

import "github.com/google/uuid"

type InstanceState string

const (
	StateRunning      InstanceState = "RUNNING"
	StateStarting                   = "STARTING"
	StateUnhealthy                  = "UNHEALTHY"
	StateShuttingDown               = "SHUTTING_DOWN"
	StateShutDown                   = "SHUT_DOWN"
)

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
