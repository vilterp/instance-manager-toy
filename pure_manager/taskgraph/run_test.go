package taskgraph

import (
	"fmt"
	"testing"

	"github.com/cockroachlabs/instance_manager/pure_manager/actions"
	"github.com/cockroachlabs/instance_manager/pure_manager/base"
)

func Test_Par(t *testing.T) {
	s := NewSpec()
	s.Par([]actions.Action{
		actions.StartInstance{Spec: base.InstanceSpec{Version: 1}},
		actions.StartInstance{Spec: base.InstanceSpec{Version: 2}},
		actions.StartInstance{Spec: base.InstanceSpec{Version: 3}},
	})
	db := NewMockGraphDB(s)
	actionRunner := actions.NewMockRunner()
	graphRunner := NewGraphRunner(db, actionRunner)
	graphRunner.Run()
	log := actionRunner.Log
	if len(log) != 4 {
		t.Fatal("need to see 4 actions")
	}
	dn := actions.DoNothing{}
	if log[3] != dn {
		t.Fatal("last one must be DoNothing")
	}
	// TODO: test that all three StartInstances show up
	fmt.Println(log)
}

func Test_Ser(t *testing.T) {
	s := NewSpec()
	s.Ser([]actions.Action{
		actions.StartInstance{Spec: base.InstanceSpec{Version: 1}},
		actions.StartInstance{Spec: base.InstanceSpec{Version: 2}},
		actions.StartInstance{Spec: base.InstanceSpec{Version: 3}},
	})
	db := NewMockGraphDB(s)
	actionRunner := actions.NewMockRunner()
	graphRunner := NewGraphRunner(db, actionRunner)
	graphRunner.Run()
	log := actionRunner.Log
	if len(log) != 3 {
		t.Fatal("need to see 3 actions")
	}
	// TODO: test that all three StartInstances show up
	fmt.Println(log)
}

func Test_Both(t *testing.T) {
	s := NewSpec()
	t1 := s.Ser([]actions.Action{
		actions.StartInstance{Spec: base.InstanceSpec{Version: 1}},
		actions.StartInstance{Spec: base.InstanceSpec{Version: 2}},
		actions.StartInstance{Spec: base.InstanceSpec{Version: 3}},
	})
	t2 := s.Ser([]actions.Action{
		actions.StartInstance{Spec: base.InstanceSpec{Version: 4}},
		actions.StartInstance{Spec: base.InstanceSpec{Version: 5}},
		actions.StartInstance{Spec: base.InstanceSpec{Version: 6}},
	})
	s.ParIDs([]TaskID{t1, t2})
	db := NewMockGraphDB(s)
	actionRunner := actions.NewMockRunner()
	graphRunner := NewGraphRunner(db, actionRunner)
	graphRunner.Run()
	log := actionRunner.Log
	fmt.Println(log)
	if len(log) != 7 {
		t.Fatal("need to see 6 actions")
	}
	// TODO: test that all three StartInstances show up
	dn := actions.DoNothing{}
	if log[6] != dn {
		t.Fatal("last one must be DoNothing")
	}
}
