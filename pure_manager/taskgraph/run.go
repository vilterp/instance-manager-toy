package taskgraph

import (
	"fmt"

	"github.com/cockroachlabs/instance_manager/pure_manager/actions"
)

type GraphRunner struct {
	db           StateDB
	actionRunner actions.Runner
	events       chan TaskEvent

	running int
	toDo    int
}

func NewGraphRunner(db StateDB, runner actions.Runner) *GraphRunner {
	return &GraphRunner{
		events:       make(chan TaskEvent),
		actionRunner: runner,
		db:           db,
		toDo:         len(db.List()),
	}
}

func (g *GraphRunner) Run() {
	g.runNext()
	for g.toDo > 0 {
		fmt.Println("todo", g.toDo, "running", g.running)
		g.runNext()
		evt := <-g.events
		switch tEvt := evt.(type) {
		case *TaskSucceeeded:
			g.db.MarkSucceeded(tEvt.ID)
			g.toDo--
			g.running--
		case *TaskFailed:
			g.db.MarkFailed(tEvt.ID, tEvt.Err)
			g.toDo--
			g.running--
		}
	}
}

func (g *GraphRunner) runNext() {
	unblocked := g.db.GetUnblockedTasks()
	if len(unblocked) == 0 && g.running == 0 {
		panic("no unblocked tasks and nothing running")
	}
	for _, t := range unblocked {
		// Go gotcha: bind vars here or they won't work in the closure
		tID := t.ID
		tAction := t.Action
		g.db.MarkStarted(tID)
		g.running++
		go func() {
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
