package ksana

import (
	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
)

type EngineHandler func(en Engine) error

type Engine interface {
	Mount(*gin.Engine)
	Seed() error
	Migrate() error
	Worker()
	Cron() map[string]func()
	Deploy()
	Shell() []cli.Command
}
