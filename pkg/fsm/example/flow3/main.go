package main

import (
	"fmt"

	"github.com/yametech/echoer/pkg/fsm"
)

func StateList(states ...string) []string { return states }

/*
flow_run myflow
	step A => (Yes -> D) {
		action = "ci";
		args = ( project="https://artifactory.compass.ym", version=12343 );
	};
	deci D => ( Yes -> B | Other -> C | No -> A ) {
		action="approval";
		args= (work_order="nz00001",version=12343);
	};
	step B => (Yes->C) {
		action="deploy";
		args=(env="release",version=12343);
	};
	step C => () {
		action="notify";
		args=(work_order="nz00001",version=12343);
	};
flow_run_end
*/

/*
[state,op,dst]
A Yes D
D Yes B
D Other C
D No A
*/

func upExample() {
	f := fsm.NewFSM(fsm.READY, nil, nil)

	f.Add(fsm.OpStart, StateList(fsm.READY), "A", func(e *fsm.Event) {
		fmt.Println("start to state A")
	})

	f.Add("A_Yes", StateList("A"), "D", func(e *fsm.Event) {
		fmt.Println("Step A accept Yes to D")
	})

	f.Add("D_Yes", StateList("D"), "B", func(e *fsm.Event) {
		fmt.Println("Step D accept Yes to B")
	})

	f.Add("D_Other", StateList("D"), "C", func(e *fsm.Event) {
		fmt.Println("Step D accept Other to B")
	})

	f.Add("D_No", StateList("D"), "A", func(e *fsm.Event) {
		fmt.Println("Step D accept No to A")
	})

	f.Add("B_Yes", StateList("B"), "C", func(e *fsm.Event) {
		fmt.Println("Step B accept Yes to C")
	})

	f.Add("C_Yes", StateList("C"), fsm.STOPPED, func(e *fsm.Event) {
		fmt.Println("Step C accept Yes to Stopped")
	})

	f.Add(fsm.OpPause, StateList("A", "B", "C", "D"), fsm.SUSPEND, func(e *fsm.Event) {
		fmt.Println("pause to state suspend")
	})

	f.Add(fsm.OpContinue, StateList(fsm.SUSPEND), "D", func(e *fsm.Event) {
		fmt.Printf("continue to state %s\n", e.Last())
	})

	f.Add(fsm.OpStop, StateList(), fsm.STOPPED, func(e *fsm.Event) {
		fmt.Println("stop to state stopped")
	})

	f.Add(fsm.READY, StateList(), "", func(e *fsm.Event) {
		fmt.Printf("my to state %s\n", e.Current())
	})

	var err error
	if err = f.Event(fsm.OpStart); err != nil {
		fmt.Println(err)
	}

	if err = f.Event("A_Yes"); err != nil {
		fmt.Println(err)
	}

	if err = f.Event(fsm.OpPause); err != nil {
		fmt.Println(err)
	}

	if err = f.Event(fsm.OpContinue); err != nil {
		fmt.Println(err)
	}

	fmt.Println("AvailableTransitions1: ", f.AvailableTransitions())

	if err = f.Event("D_Yes"); err != nil {
		fmt.Println(err)
	}

	if err = f.Event("B_Yes"); err != nil {
		fmt.Println(err)
	}

	if err = f.Event("C_Yes"); err != nil {
		fmt.Println(err)
	}
}

func main() {
	/*
		start to state A
		Step A accept Yes to D
		pause to state suspend
		continue to state suspend
		AvailableTransitions1:  [D_No D_Yes D_Other pause]
		Step D accept Yes to B
		Step B accept Yes to C
		Step C accept Yes to Stopped
	*/
	upExample()
}
