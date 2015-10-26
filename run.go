package ksana

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

var KSANA_ENV = cli.StringFlag{
	Name:   "environment, e",
	Value:  "development",
	Usage:  "can be production, development, stage, test etc...",
	EnvVar: "KSANA_ENV",
}

func Load(c *cli.Context) (*Configuration, error) {

	var cfg Configuration
	env := c.String("environment")
	if err := cfg.Load(fmt.Sprintf("config/%s/settings.toml", env)); err == nil {
		cfg.Env = env
		return &cfg, nil
	} else {
		return nil, err
	}
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
