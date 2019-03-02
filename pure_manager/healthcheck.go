package pure_manager

type HealthCheckResult interface {
	HealthCheckRes()
}

type HealthOk struct {
}

func (ho HealthOk) HealthCheckRes() {}

type HealthErr struct {
	err error
}

func (he HealthErr) HealthCheckRes() {}

type HealthChecker interface {
	HealthCheck(id InstanceID) HealthCheckResult
}
