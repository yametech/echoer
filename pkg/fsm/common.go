package fsm

import (
	"fmt"
)

type OperatorStateType = string

const (
	READY   OperatorStateType = "ready"
	RUNNING OperatorStateType = "running"
	SUSPEND OperatorStateType = "suspend"
	STOPPED OperatorStateType = "stopped"
	DONE    OperatorStateType = "done"
)

type OperatorType = string

const (
	OpStart    OperatorType = "start"
	OpStop     OperatorType = "stop"
	OpPause    OperatorType = "pause"
	OpContinue OperatorType = "continue"
	OpEnd      OperatorType = "end"
)

type EventTriggerMechanismType = string // ["before,after,enter,leave"]

const (
	BEFORE EventTriggerMechanismType = "before"
	AFTER  EventTriggerMechanismType = "after"
	ENTER  EventTriggerMechanismType = "enter"
	LEAVE  EventTriggerMechanismType = "leave"
)

func EventTriggerMechanismTypePrefix(et EventTriggerMechanismType) string {
	return fmt.Sprintf("%s_", et)
}

type EventOrStateNameType = string

const (
	EVENT EventOrStateNameType = "event"
	STATE EventOrStateNameType = "state"
)

type EventOrStateType = string // ["before_{}"]

const (
	BeforeEvent = BEFORE + "_" + EVENT
	LeaveState  = LEAVE + "_" + STATE
	EnterState  = ENTER + "_" + STATE
	AfterEvent  = AFTER + "_" + EVENT
)

func BeforeEventOrState(eventOrStateName EventOrStateNameType) EventOrStateType {
	return fmt.Sprintf("%s_%s", BEFORE, eventOrStateName)
}

func AfterEventOrState(eventOrStateName EventOrStateNameType) EventOrStateType {
	return fmt.Sprintf("%s_%s", AFTER, eventOrStateName)
}

func EnterEventOrState(eventOrStateName EventOrStateNameType) EventOrStateType {
	return fmt.Sprintf("%s_%s", ENTER, eventOrStateName)
}

func LeaveEventOrState(eventOrStateName EventOrStateNameType) EventOrStateType {
	return fmt.Sprintf("%s_%s", LEAVE, eventOrStateName)
}
