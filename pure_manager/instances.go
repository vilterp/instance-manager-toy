package pure_manager

const (
	StateRunning      = "RUNNING"
	StateStarting     = "STARTING"
	StateUnhealthy    = "UNHEALTHY"
	StateShuttingDown = "SHUTTING_DOWN"
	StateShutDown     = "SHUT_DOWN"
)

type InstanceStateDB interface {
	List() []*Instance
	ListHealthy() []*Instance

	UpdateHealthStatus(id InstanceID, res HealthCheckResult)
}

type InstanceID int

type Version int

type Instance struct {
	ID      InstanceID
	Version Version
}

type InstanceSpec struct {
	Version Version
}

// TODO: runner
