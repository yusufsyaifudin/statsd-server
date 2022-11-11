package main

import (
	"fmt"
	"github.com/mitchellh/cli"
	statsd_server "github.com/yusufsyaifudin/statsd-server/cmd/statsd-server"
	"log"
	"os"
)

const (
	CliName    = "statsd-server"
	CliVersion = "1.0.0"
)

func main() {
	serverCmdCfg := statsd_server.Config{}

	c := cli.NewCLI(CliName, CliVersion)
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"": func() (cli.Command, error) {
			return statsd_server.New(serverCmdCfg)
		},
		"server": func() (cli.Command, error) {
			return statsd_server.New(serverCmdCfg)
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		err = fmt.Errorf("error main program: %w", err)
		log.Println(err)
	}

	os.Exit(exitStatus)
}
