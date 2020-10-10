package command

import "fmt"

//ErrWrongTypeOp means that operation is not acceptable for the given key.
var ErrWrongTypeOp = fmt.Errorf("command: wrong type operation")

type Reply interface {
	Value() interface{}
}

type Command interface {
	//Name returns the command name.
	Name() string
	//Help returns information about the command. Description, usage and etc.
	Help() string
	//Execute executes the command with the given arguments.
	Execute(args ...string) Reply
}

type CommandParser interface {
	Parse(str string) (cmd Command, args []string, err error)
}
