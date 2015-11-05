package orm

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/itpkg/ksana/cmd"
	"github.com/itpkg/ksana/logging"
)

func init() {
	cmd.Register(cli.Command{
		Name:    "database",
		Aliases: []string{"db"},
		Usage:   "database operations",
		Subcommands: []cli.Command{
			{
				Name:    "seed",
				Aliases: []string{"s"},
				Usage:   "load the seed data into database",
				Flags:   []cli.Flag{cmd.FLAG_ENV},
				Action: func(c *cli.Context) {
					//todo
					env := cmd.Env(c)
					log := logging.Open(env)
					log.Debug("env = [%s, %s]", env, os.Getenv("KSANA_ENV"))

				},
			},
			{
				Name:    "migrate",
				Aliases: []string{"m"},
				Usage:   "migrate the database",
				Action: func(c *cli.Context) {
					//todo
				},
			},
		},
	})
}
