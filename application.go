package ksana

import (
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/jrallison/go-workers"
	"github.com/op/go-logging"
)

type Application struct {
	Cfg    *Configuration  `inject:""`
	Logger *logging.Logger `inject:""`
	Router *gin.Engine     `inject:""`
	Redis  *redis.Pool     `inject:""`
	Db     *gorm.DB        `inject:""`
}

//==============================================================================

func New(c *cli.Context) (*Application, error) {
	cfg, err := Load(c)
	if err != nil {
		return nil, err
	}
	var db *gorm.DB
	if db, err = cfg.Db(); err != nil {
		return nil, err
	}

	app := Application{}
	if err = Use(&app, db, cfg); err != nil {
		return nil, err
	}
	if err = beans.Populate(); err != nil {
		return nil, err
	}

	//------workers---------------------------
	workers.Configure(map[string]string{
		"server":   cfg.Redis.Server(),
		"database": strconv.Itoa(cfg.Redis.Db),
		"pool":     "10",
		"process":  "1",
	})

	return &app, nil
}
