package ksana

import (
	"github.com/codegangsta/cli"
)

type EngineHandler func(en Engine) error

type Engine interface {
	Router()
	Migrate() error
	Job()
	Deploy()
	Shell() []cli.Command
}
