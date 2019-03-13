package taskgraph

import (
	"github.com/vilterp/instance-manager-toy/actions"
	"github.com/vilterp/instance-manager-toy/db"
	"github.com/vilterp/instance-manager-toy/proto"
)

type Runner struct {
	db           db.TasksDB
	actionRunner actions.Runner
	events       chan *proto.TaskEvent

	running int
	toDo    int
}

func NewRunner(db db.TasksDB, runner actions.Runner) *Runner {
	return &Runner{
		events:       make(chan *proto.TaskEvent),
		actionRunner: runner,
		db:           db,
		toDo:         len(db.List()),
	}
}

func (g *Runner) Run() {
	g.runNext()
	for g.toDo > 0 {
		g.runNext()
		evt := <-g.events
		switch tEvt := evt.Event.(type) {
		case *proto.TaskEvent_Succeeded_:
			succ := tEvt.Succeeded
			g.db.MarkSucceeded(proto.TaskID(succ.Id))
			g.toDo--
			g.running--
		case *proto.TaskEvent_Failed_:
			fail := tEvt.Failed
			g.db.MarkFailed(proto.TaskID(fail.Id), fail.Error)
			g.toDo--
			g.running--
		}
	}
}

func (g *Runner) runNext() {
	unblocked := g.db.GetUnblockedTasks()
	if len(unblocked) == 0 && g.running == 0 {
		panic("no unblocked tasks and nothing running")
	}
	for _, t := range unblocked {
		// Go gotcha: bind vars here or they won't work in the closure
		tID := proto.TaskID(t.Id)
		tAction := t.Action
		g.db.MarkStarted(tID)
		g.running++
		go func() {
			err := g.actionRunner.Run(tAction)
			if err != nil {
				// Jesus, these are ridiculously verbose
				g.events <- &proto.TaskEvent{
					Event: &proto.TaskEvent_Failed_{
						Failed: &proto.TaskEvent_Failed{
							Error: err.Error(),
							Id:    string(tID),
						},
					},
				}
			} else {
				g.events <- &proto.TaskEvent{
					Event: &proto.TaskEvent_Succeeded_{
						Succeeded: &proto.TaskEvent_Succeeded{
							Id: string(tID),
						},
					},
				}
			}
		}()
	}
}
