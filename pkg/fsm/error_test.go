package fsm

import (
	"errors"
	"testing"
)

func TestInvalidEventError(t *testing.T) {
	event := "invalid event"
	state := "state"
	e := InvalidEventError{Event: event, State: state}
	if e.Error() != "event "+e.Event+" inappropriate in current state "+e.State {
		t.Error("InvalidEventError string mismatch")
	}
}

func TestUnknownEventError(t *testing.T) {
	event := "invalid event"
	e := UnknownEventError{Event: event}
	if e.Error() != "event "+e.Event+" does not exist" {
		t.Error("UnknownEventError string mismatch")
	}
}

func TestInTransitionError(t *testing.T) {
	event := "in transition"
	e := InTransitionError{Event: event}
	if e.Error() != "event "+e.Event+" inappropriate because previous transition did not complete" {
		t.Error("InTransitionError string mismatch")
	}
}

func TestNotInTransitionError(t *testing.T) {
	e := NotInTransitionError{}
	if e.Error() != "transition inappropriate because no state change in progress" {
		t.Error("NotInTransitionError string mismatch")
	}
}

func TestNoTransitionError(t *testing.T) {
	e := NoTransitionError{}
	if e.Error() != "no transition" {
		t.Error("NoTransitionError string mismatch")
	}
	e.Err = errors.New("no transition")
	if e.Error() != "no transition with error: "+e.Err.Error() {
		t.Error("NoTransitionError string mismatch")
	}
}

func TestCanceledError(t *testing.T) {
	e := CanceledError{}
	if e.Error() != "transition canceled" {
		t.Error("CanceledError string mismatch")
	}
	e.Err = errors.New("canceled")
	if e.Error() != "transition canceled with error: "+e.Err.Error() {
		t.Error("CanceledError string mismatch")
	}
}

func TestAsyncError(t *testing.T) {
	e := AsyncError{}
	if e.Error() != "async started" {
		t.Error("AsyncError string mismatch")
	}
	e.Err = errors.New("async")
	if e.Error() != "async started with error: "+e.Err.Error() {
		t.Error("AsyncError string mismatch")
	}
}

func TestInternalError(t *testing.T) {
	e := InternalError{}
	if e.Error() != "internal error on state transition" {
		t.Error("InternalError string mismatch")
	}
}
