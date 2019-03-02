package manager

type InstanceState string

const (
	StateRunning      = "RUNNING"
	StateStarting     = "STARTING"
	StateUnhealthy    = "UNHEALTHY"
	StateShuttingDown = "SHUTTING_DOWN"
	StateShutDown     = "SHUT_DOWN"
)

type InstanceID int

type Instance struct {
	ID      InstanceID
	State   InstanceState
	Version int
}

type Runner interface {
	Start(spec InstanceSpec) (*Instance, *Operation, error)
	ShutDown(id InstanceID) *Operation
	GetInstance(id InstanceID) *Instance
	ListInstances() []*Instance
	ListUpInstances() []*Instance
	GetOpLog() OpLogReader
}

type HealthChecker interface {
	IsHealthy(instanceID int) bool
}
