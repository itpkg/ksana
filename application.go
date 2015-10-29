package ksana

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
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
	//---------database-----------------------
	var db *gorm.DB
	if db, err = cfg.OpenDb(); err != nil {
		return nil, err
	}
	//----------encrypt-----------------------
	var a_c cipher.Block
	if a_c, err = aes.NewCipher(cfg.Secrets[90:122]); err != nil {
		return nil, err
	}

	//----------application-------------------
	app := Application{}
	if err = Use(
		&app,
		db,
		cfg,
		cfg.OpenRedis(),
	); err != nil {
		return nil, err
	}
	if err = Map(map[string]interface{}{
		"aes.cipher": a_c,
		"hmac.key":   cfg.Secrets[50:82],
		"hmac.fn":    sha512.New,
	}); err != nil {
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
