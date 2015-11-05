package orm

import (
	"log"

	"github.com/codegangsta/cli"
	"github.com/itpkg/ksana/cmd"
	"github.com/itpkg/ksana/logging"
	"github.com/itpkg/ksana/utils"
)

func cli_cfg(fn func(*Configuration) error) func(*cli.Context) {
	return func(c *cli.Context) {
		env := cmd.Env(c)
		cfg := Configuration{}
		err := cfg.Load(env)
		if err == nil {
			err = fn(&cfg)
		}
		if err == nil {
			log.Println("Done.")
		} else {
			log.Fatal(err)
		}
	}
}

func cli_db(fn func(*Db) error) func(*cli.Context) {
	return func(c *cli.Context) {
		env := cmd.Env(c)
		log := logging.Open(env)
		db, err := Open(env)
		if err == nil {
			db.Logger = log
			err = fn(db)
		}

		if err == nil {
			log.Info("Done!")
		} else {
			log.Error(err.Error())
		}

	}
}

func init() {
	cmd.Register(cli.Command{
		Name:    "database",
		Aliases: []string{"db"},
		Usage:   "database operations",
		Flags:   []cli.Flag{cmd.FLAG_ENV},
		Subcommands: []cli.Command{
			{
				Name:    "console",
				Aliases: []string{"c"},
				Usage:   "database console",
				Flags:   []cli.Flag{cmd.FLAG_ENV},
				Action: cli_cfg(func(cfg *Configuration) error {
					c, a := cfg.Connect()
					return utils.Shell(c, a...)
				}),
			},
			{
				Name:    "create",
				Aliases: []string{"n"},
				Usage:   "create database",
				Flags:   []cli.Flag{cmd.FLAG_ENV},
				Action: cli_cfg(func(cfg *Configuration) error {
					c, a := cfg.Create()
					return utils.Shell(c, a...)
				}),
			},
			{
				Name:    "drop",
				Aliases: []string{"d"},
				Usage:   "drop database",
				Flags:   []cli.Flag{cmd.FLAG_ENV},
				Action: cli_cfg(func(cfg *Configuration) error {
					c, a := cfg.Drop()
					return utils.Shell(c, a...)
				}),
			},
			{
				Name:    "seed",
				Aliases: []string{"s"},
				Usage:   "load the seed data into database",
				Flags:   []cli.Flag{cmd.FLAG_ENV},
				Action: cli_db(func(db *Db) error {
					//todo
					return nil
				}),
			},
			{
				Name:    "migrate",
				Aliases: []string{"m"},
				Usage:   "migrate the database",
				Flags:   []cli.Flag{cmd.FLAG_ENV},
				Action: cli_db(func(db *Db) error {
					return db.Migrate()
				}),
			},
			{
				Name:    "rollback",
				Aliases: []string{"r"},
				Usage:   "rollback the database",
				Flags:   []cli.Flag{cmd.FLAG_ENV},
				Action: cli_db(func(db *Db) error {
					return db.Rollback()
				}),
			},
			{
				Name:    "backup",
				Aliases: []string{"b"},
				Usage:   "backup the database",
				Flags:   []cli.Flag{cmd.FLAG_ENV},
				Action: cli_cfg(func(cfg *Configuration) error {
					c, a := cfg.Backup()
					return utils.Shell(c, a...)
				}),
			},
		},
	})
}
