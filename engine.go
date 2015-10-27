package ksana

import (
	"github.com/codegangsta/cli"
)

type EngineHandler func(en Engine) error

type Engine interface {
	Router()
	Seed() error
	Migrate() error
	Job()
	Deploy()
	Shell() []cli.Command
}
