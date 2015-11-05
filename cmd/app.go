package cmd

import (
	"os"

	"github.com/codegangsta/cli"
)

var commands = make([]cli.Command, 0)

func Register(args ...cli.Command) {
	commands = append(commands, args...)
}

func Run() error {
	app := cli.NewApp()
	app.Name = "ksana"
	app.Usage = "A web framwork for go like rails"
	app.Version = "v20151104"
	app.Commands = commands

	return app.Run(os.Args)
}
