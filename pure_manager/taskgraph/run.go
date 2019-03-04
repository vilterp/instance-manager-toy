package taskgraph

import (
	"fmt"
	"log"

	"github.com/cockroachlabs/instance_manager/pure_manager/actions"
)

type GraphRunner struct {
	db           StateDB
	actionRunner actions.Runner
	events       chan TaskEvent
}

func NewGraphRunner(db StateDB, runner actions.Runner) *GraphRunner {
	return &GraphRunner{
		events:       make(chan TaskEvent),
		actionRunner: runner,
		db:           db,
	}
}

func (g *GraphRunner) Run() {
	toDo := len(g.db.List())
	running := 0

	g.runNext(running)
	for toDo > 0 {
		evt := <-g.events
		switch tEvt := evt.(type) {
		case *TaskStarted:
			g.db.MarkStarted(tEvt.ID)
			running++
		case *TaskSucceeeded:
			g.db.MarkSucceeded(tEvt.ID)
			toDo--
			fmt.Println("succeeded; run next")
			g.runNext(running)
		case *TaskFailed:
			g.db.MarkFailed(tEvt.ID, tEvt.Err)
			toDo--
			g.runNext(running)
		}
	}
}

func (g GraphRunner) runNext(numRunning int) {
	unblocked := g.db.GetUnblockedTasks()
	if len(unblocked) == 0 && numRunning == 0 {
		panic("no unblocked tasks and nothing running")
	}
	log.Println("unblocked:", unblocked)
	for _, t := range unblocked {
		// Go gotcha: bind vars here or they won't work in the closure
		tID := t.ID
		tAction := t.Action
		fmt.Println("about to run", tID.String())
		go func() {
			g.events <- &TaskStarted{tID}
			err := g.actionRunner.Run(tAction)
			if err != nil {
				g.events <- &TaskFailed{
					Err: err,
					ID:  tID,
				}
			} else {
				g.events <- &TaskSucceeeded{tID}
			}
		}()
	}
}

type TaskEvent interface {
	TaskID() TaskID
}

type TaskStarted struct {
	ID TaskID
}

func (ts *TaskStarted) TaskID() TaskID {
	return ts.ID
}

type TaskSucceeeded struct {
	ID TaskID
}

func (ts *TaskSucceeeded) TaskID() TaskID {
	return ts.ID
}

type TaskFailed struct {
	ID  TaskID
	Err error
}

func (tf *TaskFailed) TaskID() TaskID {
	return tf.ID
}
