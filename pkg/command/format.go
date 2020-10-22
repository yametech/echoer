package command

import (
	tw "github.com/olekukonko/tablewriter"
	"io"
)

var (
	_ Formatter = &Format{}
	_ io.Writer = &Format{}
)

type Formatter interface {
	Header(...string) Formatter
	Row(...string) Formatter
	Out() []byte
}

type Format struct {
	table *tw.Table
	self  []byte
}

func NewFormat() *Format {
	format := &Format{
		self: make([]byte, 0),
	}
	table := tw.NewWriter(format)
	format.table = table
	return format
}

func (f *Format) Write(p []byte) (n int, err error) {
	f.self = append(f.self, p...)
	return len(p), nil
}

func (f *Format) Header(s ...string) Formatter {
	f.table.SetHeader(s)
	return f
}

func (f *Format) Row(s ...string) Formatter {
	f.table.Append(s)
	return f
}

func (f *Format) Out() []byte {
	f.table.Render()
	return f.self
}
