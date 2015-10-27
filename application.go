package ksana

import (
	"github.com/codegangsta/cli"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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
	app := Application{}
	if err = Use(&app, cfg); err != nil {
		return nil, err
	}
	if err = beans.Populate(); err != nil {
		return nil, err
	}

	return &app, nil
}
