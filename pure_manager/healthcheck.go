package pure_manager

import "github.com/cockroachlabs/instance_manager/pure_manager/base"

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
	HealthCheck(id base.InstanceID) HealthCheckResult
}