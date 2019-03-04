package actions

import "log"

type Runner interface {
	Run(a Action) error
}

type MockRunner struct {
	log []Action
}

var _ Runner = &MockRunner{}

func (m MockRunner) Run(a Action) error {
	log.Println("running", a.String())
	m.log = append(m.log, a)
	return nil
}
