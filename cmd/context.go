package cmd

import (
	"github.com/codegangsta/cli"
)

var FLAG_ENV = cli.StringFlag{
	Name:   "environment, e",
	Value:  "development",
	Usage:  "can be production, development, stage, test etc...",
	EnvVar: "KSANA_ENV",
}

func Env(c *cli.Context) string {
	return c.String("environment")
}

func IsProduction(env string) bool {
	return env == "production"
}
