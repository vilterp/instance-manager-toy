package taskgraph

func NewSpec() *Spec {
	return &Spec{
		Tasks: map[TaskID]Spec{},
	}
}

func (s *Spec) Par(ts []string) {

}
