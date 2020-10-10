package main

import (
	"fmt"

	"github.com/yametech/echoer/pkg/fsm"
)

type IAction interface {
	Handle(*fsm.Event)
}

type AbstractAction struct{}

func (a *AbstractAction) Do(event *fsm.Event) error {
	panic("not implement IAction Do")
}

var _ IAction = &nFoundAction{}
var _ IAction = &iFoundAction{}
var _ IAction = &cFoundAction{}

type MyAction1Params = map[string]interface{}

type nFoundAction struct {
	AbstractAction
}

func (ma *nFoundAction) Handle(event *fsm.Event) {
	fmt.Printf("nFoundAction event = %s\n", event.Event)
	fmt.Printf("flow current state=%s\n", event.Current())
}

type iFoundAction struct {
	AbstractAction
}

func (ma *iFoundAction) Handle(event *fsm.Event) {
	fmt.Printf("iFoundAction event = %s\n", event.Event)
	fmt.Printf("flow current state=%s\n", event.Current())
}

type cFoundAction struct {
	AbstractAction
}

func (ma *cFoundAction) Handle(event *fsm.Event) {
	fmt.Printf("cFoundAction event = %s\n", event.Event)
	fmt.Printf("flow current state=%s\n", event.Current())
}

type Flow2 struct {
	Name string `json:"name"`
	*fsm.FSM
}

func NewFlow2(name string) *Flow2 {
	flow2 := &Flow2{
		Name: name,
		FSM:  fsm.NewFSM(fsm.READY, nil, nil),
	}
	return flow2
}

func main() {

	flow2 := NewFlow2("flow2-example")
	flow2.Add(fsm.OpStart, []string{fsm.READY}, fsm.RUNNING, nil)
	flow2.Add("n", []string{fsm.RUNNING}, "n_found", (&nFoundAction{}).Handle)
	flow2.Add("i", []string{"n_found"}, "i_found", (&iFoundAction{}).Handle)
	flow2.Add("c", []string{"i_found"}, "c_found", (&cFoundAction{}).Handle)
	flow2.Add(fsm.OpStop, []string{"c_found"}, fsm.STOPPED, nil)

	if err := flow2.Send(fsm.OpStart); err != nil {
		panic(err)
	}

	if err := flow2.Send("n"); err != nil {
		panic(err)
	}

	if err := flow2.Send("i"); err != nil {
		panic(err)
	}

	if err := flow2.Send("c"); err != nil {
		panic(err)
	}

	if err := flow2.Send(fsm.OpStop); err != nil {
		panic(err)
	}
	fmt.Printf("flow2 current state=%s\n", flow2.Current())
}
