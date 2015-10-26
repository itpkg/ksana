package ksana

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

func Run() error {
	app := cli.NewApp()
	app.Name = "ksana"
	app.Usage = "A web framwork for go like rails"
	app.Version = "v20151026"

	commands := []cli.Command{}
	if err := LoopEngine(func(en Engine) error {
		cds := en.Shell()
		commands = append(commands, cds...)
		return nil
	}); err != nil {
		return err
	}
	app.Commands = commands

	return app.Run(os.Args)
}
