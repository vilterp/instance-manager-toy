package taskgraph

import "github.com/cockroachlabs/instance_manager/pure_manager/actions"

type failureReport struct {
	taskID TaskID
	err    error
}

type GraphRunner struct {
	db           StateDB
	actionRunner actions.Runner
	started      chan TaskID
	succeeded    chan TaskID
	failed       chan failureReport
}

func NewGraphRunner(db StateDB, runner actions.Runner) *GraphRunner {
	return &GraphRunner{
		started:      make(chan TaskID),
		succeeded:    make(chan TaskID),
		failed:       make(chan failureReport),
		actionRunner: runner,
		db:           db,
	}
}

func (g *GraphRunner) Run() {
	toDo := len(g.db.List())

	for toDo > 0 {
		select {
		case succID := <-g.succeeded:
			g.db.MarkSucceeded(succID)
			toDo--
			g.runNext()
		case report := <-g.failed:
			g.db.MarkFailed(report.taskID, report.err)
			toDo--
			g.runNext()
		}
	}
}

func (g GraphRunner) runNext() {
	unblocked := g.db.GetUnblockedTasks()
	for _, t := range unblocked {
		go func() {
			g.started <- t.ID
			err := g.actionRunner.Run(t.Action)
			if err != nil {
				g.failed <- failureReport{
					err:    err,
					taskID: t.ID,
				}
			} else {
				g.succeeded <- t.ID
			}
		}()
	}
}
