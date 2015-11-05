package ioc

import (
	"os"

	"github.com/codegangsta/cli"
)

var KSANA_ENV = cli.StringFlag{
	Name:   "environment, e",
	Value:  "development",
	Usage:  "can be production, development, stage, test etc...",
	EnvVar: "KSANA_ENV",
}

var app = cli.NewApp()

func Run() error {
	return app.Run(os.Args)
}

func Command(cmds ...cli.Command) {
	app.Commands = append(app.Commands, cmds...)
}

func init() {
	app.Name = "ksana"
	app.Usage = "A web framwork for go like rails"
	app.Version = "v20151104"
	app.Commands = make([]cli.Command, 0)
}
