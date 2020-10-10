package command

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"

	"github.com/yametech/echoer/pkg/storage"
)

//ErrCommandNotFound means that command could not be parsed.
var ErrCommandNotFound = fmt.Errorf("command: not found")

//Parser is a parser that parses user input and creates the appropriate command.
type Parser struct {
	storage.IStorage
}

//NewParser creates a new parser
func NewParser(storage storage.IStorage) *Parser {
	return &Parser{storage}
}

//Parse parses string to Command with args
func (p *Parser) Parse(str string) (Command, []string, error) {
	var cmd Command
	trimPrefixStr := strings.TrimSpace(str)
	switch {
	case strings.HasPrefix(strings.ToLower(trimPrefixStr), "flow_run"):
		cmd = &FlowRunCmd{[]byte(str), p.IStorage}
		return cmd, nil, nil
	case strings.HasPrefix(strings.ToLower(trimPrefixStr), "flow"):
		cmd = &FlowCmd{[]byte(str), p.IStorage}
		return cmd, nil, nil
	case strings.HasPrefix(strings.ToLower(trimPrefixStr), "action"):
		cmd = &ActionCmd{[]byte(str), p.IStorage}
		return cmd, nil, nil
	}

	args := p.extractArgs(trimPrefixStr)
	if len(args) == 0 {
		return nil, nil, ErrCommandNotFound
	}

	switch strings.ToLower(args[0]) {
	case "list":
		cmd = &List{p.IStorage}
	case "get":
		cmd = &Get{p.IStorage}
	case "del":
		cmd = &Del{p.IStorage}
	case "help":
		cmd = &Help{}
	default:
		return nil, nil, ErrCommandNotFound
	}

	return cmd, args[1:], nil
}

func (p *Parser) extractArgs(val string) []string {
	args := make([]string, 0)
	var inQuote bool
	var buf bytes.Buffer
	for _, r := range val {
		switch {
		case r == '`':
			inQuote = !inQuote
		case unicode.IsSpace(r):
			if !inQuote && buf.Len() > 0 {
				args = append(args, buf.String())
				buf.Reset()
			} else {
				buf.WriteRune(r)
			}
		default:
			buf.WriteRune(r)
		}
	}
	if buf.Len() > 0 {
		args = append(args, buf.String())
	}
	return args
}
