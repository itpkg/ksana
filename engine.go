package ksana

import (
	"github.com/codegangsta/cli"
)

type EngineHandler func(en Engine) error

type Engine interface {
	Router()
	Seed() error
	Migrate() error
	Worker()
	Cron() map[string]func()
	Deploy()
	Shell() []cli.Command
}
