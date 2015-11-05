package cmd

import (
	"github.com/codegangsta/cli"
)

var commands = make([]cli.Command, 0)

func Register(args ...cli.Command) {
	commands = append(commands, args...)
}

func New(name, usage, version string) *cli.App {
	app := cli.NewApp()
	app.Name = name
	app.Usage = usage
	app.Version = version
	app.Commands = commands
	return app
}
