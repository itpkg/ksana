package ksana

import (
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

func (p *Application) DbMigrate() error {
	return LoopEngine(func(en Engine) error {
		return en.Migrate()
	})
}
