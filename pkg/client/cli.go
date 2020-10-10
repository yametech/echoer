package client

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/yametech/echoer/api"
	"google.golang.org/grpc"
)

const connectTimeout = 200 * time.Millisecond

const (
	prefix = "> "
)

var space = ""

//CLI allows users to interact with a server.
type CLI struct {
	printer *printer
	term    *Prompt
	conn    *grpc.ClientConn
	client  api.EchoClient
}

//Run runs a new CLI.
func Run(hostPorts string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, hostPorts, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("could not dial %s: %v", hostPorts, err)
	}
	term := NewPrompt()
	space = fmt.Sprintf("%s%s", "", prefix)
	term.SetPrefix(space)

	c := &CLI{
		printer: newPrinter(os.Stdout),
		term:    term,
		client:  api.NewEchoClient(conn),
		conn:    conn,
	}
	defer func() { _ = c.Close() }()

	c.run()

	return nil
}

func (c *CLI) Close() error {
	if err := c.printer.Close(); err != nil {
		return err
	}
	if err := c.conn.Close(); err != nil {
		return err
	}
	return nil
}

func (c *CLI) run() {
	c.printer.printLogo()
	h := func(command string) {
		req := &api.ExecuteRequest{Command: []byte(command)}
		space = fmt.Sprintf("%s%s", "", prefix)
		if resp, err := c.client.Execute(context.Background(), req); err != nil {
			c.printer.printError(err)
		} else {
			c.printer.printResponse(resp)
		}
	}
	c.term.Handler(h)
	c.printer.println("Bye!")
}
