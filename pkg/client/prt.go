package client

import (
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
	"github.com/yametech/echoer/api"
)

const (
	okString  = "OK"
	nilString = "(nil)"
)

const logo = `
echoer
        /\_/\                                                 ##         .
      =( °w° )=                                         ## ## ##        ==
        )   (     // 📒 🤔🤔🤔  ♻︎                       ## ## ## ## ##    ===
       (__ __)           === == ==                /""""""""""""""""\___/ ===
 /"""""""""""""" //\___/ === == ==                           ~~/~~~~ ~~~ ~~~~ ~~ ~ /  ===- ~~~
{                       /  == =-                  \______ o          _,/
 \______ O           _ _/                          \      \       _,'
  \    \         _ _/                               '--.._\..--''
    \____\_______/__/__/
`

type printer struct {
	okColor  *color.Color
	errColor *color.Color
	nilColor *color.Color
	out      io.Writer
}

func newPrinter(out io.Writer) *printer {
	return &printer{
		okColor:  color.New(color.FgHiGreen),
		errColor: color.New(color.FgHiRed),
		nilColor: color.New(color.FgHiCyan),
		out:      out,
	}
}

//Close closes the printer
func (p *printer) Close() error {
	if cl, ok := p.out.(io.Closer); ok {
		return cl.Close()
	}
	return nil
}

func (p *printer) printLogo() {
	color.Set(color.FgMagenta)
	p.println(strings.Replace(logo, "\n", "\r\n", -1))
	color.Unset()
}

func (p *printer) println(str string) {
	_, _ = fmt.Fprintf(p.out, "%s\r\n", str)
}

func (p *printer) printError(err error) {
	_, _ = p.errColor.Fprintf(p.out, "(ERROR): %s\n", err.Error())
}

func (p *printer) printResponse(resp *api.ExecuteCommandResponse) {
	switch resp.Reply {
	case api.CommandExecutionReply_OK:
		p.println(p.okColor.Sprint(okString))
	case api.CommandExecutionReply_NIL:
		p.println(p.nilColor.Sprint(nilString))
	case api.CommandExecutionReply_Raw:
		p.println(fmt.Sprintf("R| %s", resp.Raw))
	case api.CommandExecutionReply_ERR:
		_, _ = p.errColor.Fprintf(p.out, "E| %s\n", resp.Raw)
	default:
		_, _ = fmt.Fprintf(p.out, "%v\n", resp)
	}
}
