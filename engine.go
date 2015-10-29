package ksana

import (
	"github.com/codegangsta/cli"
	"github.com/gorilla/mux"
)

type EngineHandler func(en Engine) error

type Engine interface {
	Mount(Router *mux.Router)
	Seed() error
	Migrate() error
	Worker()
	Cron() map[string]func()
	Deploy()
	Shell() []cli.Command
}
