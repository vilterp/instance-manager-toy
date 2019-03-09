package taskgraph

import (
	"log"

	"github.com/cockroachlabs/instance_manager/pure_manager/actions"
	"github.com/cockroachlabs/instance_manager/pure_manager/db"
	"github.com/cockroachlabs/instance_manager/pure_manager/proto"
)

type GraphRunner struct {
	db           db.TasksDB
	actionRunner actions.Runner
	events       chan *proto.TaskEvent

	running int
	toDo    int
}

func NewGraphRunner(db db.TasksDB, runner actions.Runner) *GraphRunner {
	return &GraphRunner{
		events:       make(chan *proto.TaskEvent),
		actionRunner: runner,
		db:           db,
		toDo:         len(db.List()),
	}
}

func (g *GraphRunner) Run() {
	g.runNext()
	for g.toDo > 0 {
		log.Println("todo", g.toDo, "running", g.running)
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
	log.Println("finished graph")
}

func (g *GraphRunner) runNext() {
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
