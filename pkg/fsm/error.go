package fsm

// InvalidEventError is returned by FSM.Event() when the event cannot be called
// in the current state.
type InvalidEventError struct {
	Event string
	State string
}

func (e InvalidEventError) Error() string {
	return "event " + e.Event + " inappropriate in current state " + e.State
}

// UnknownEventError is returned by FSM.Event() when the event is not defined.
type UnknownEventError struct {
	Event string
}

func (e UnknownEventError) Error() string {
	return "event " + e.Event + " does not exist"
}

// InTransitionError is returned by FSM.Event() when an asynchronous transition
// is already in progress.
type InTransitionError struct {
	Event string
}

func (e InTransitionError) Error() string {
	return "event " + e.Event + " inappropriate because previous transition did not complete"
}

// NotInTransitionError is returned by FSM.Transition() when an asynchronous
// transition is not in progress.
type NotInTransitionError struct{}

func (e NotInTransitionError) Error() string {
	return "transition inappropriate because no state change in progress"
}

// NoTransitionError is returned by FSM.Event() when no transition have happened,
// for example if the source and destination states are the same.
type NoTransitionError struct {
	Err error
}

func (e NoTransitionError) Error() string {
	if e.Err != nil {
		return "no transition with error: " + e.Err.Error()
	}
	return "no transition"
}

// CanceledError is returned by FSM.Event() when a callback have canceled a
// transition.
type CanceledError struct {
	Err error
}

func (e CanceledError) Error() string {
	if e.Err != nil {
		return "transition canceled with error: " + e.Err.Error()
	}
	return "transition canceled"
}

// AsyncError is returned by FSM.Event() when a callback have initiated an
// asynchronous state transition.
type AsyncError struct {
	Err error
}

func (e AsyncError) Error() string {
	if e.Err != nil {
		return "async started with error: " + e.Err.Error()
	}
	return "async started"
}

// InternalError is returned by FSM.Event() and should never occur. It is a
// probably because of a bug.
type InternalError struct{}

func (e InternalError) Error() string {
	return "internal error on state transition"
}
