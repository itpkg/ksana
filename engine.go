package ksana

import (
	"github.com/codegangsta/cli"
	"github.com/jinzhu/gorm"
)

type EngineHandler func(en Engine) error

type Engine interface {
	Router()
	Migrate(*gorm.DB) error
	Job()
	Deploy()
	Shell() []cli.Command
}
