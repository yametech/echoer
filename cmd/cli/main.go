package main

import (
	"flag"
	"fmt"
	"os"

	cli "github.com/yametech/echoer/pkg/client"
)

var (
	version = "v0.0.1"
	commit  = "v0.0.10"
	date    = "20200923"
)

func main() {
	var hosts = flag.String("host", "127.0.0.1:8081", "Host to connect to a server.")
	var showVersion = flag.Bool("version", false, "Show source-raw version.")

	flag.Parse()

	if *showVersion {
		fmt.Printf("version: %s\ncommit: %s\nbuildtime: %s", version, commit, date)
		os.Exit(0)
	}

	if err := cli.Run(*hosts); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "could not run CLI: %v", err)
	}
}
