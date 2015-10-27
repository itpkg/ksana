package ksana

import (
	"github.com/codegangsta/cli"
	"github.com/jinzhu/gorm"
)

type EngineHandler func(en Engine) error

type Engine interface {
	Router()
	Seed(*gorm.DB) error
	Migrate(*gorm.DB) error
	Job()
	Deploy()
	Shell() []cli.Command
}
