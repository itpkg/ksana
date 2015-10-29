package ksana

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/codegangsta/cli"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
)

type Configuration struct {
	Env           string           `toml:"-"`
	Secrets       []byte           `toml:"-"`
	SecretsS      string           `toml:"secrets"`
	Http          HttpCfg          `toml:"http"`
	Database      DatabaseCfg      `toml:"database"`
	Redis         RedisCfg         `toml:"redis"`
	Elasticsearch ElasticsearchCfg `toml:"elasticsearch"`
	Workers       map[string]int   `toml:"workers"`
}

type HttpCfg struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type DatabaseCfg struct {
	Dialect string  `toml:"dialect"`
	Url     string  `toml:"url"`
	Pool    PoolCfg `toml:"pool"`
}

type RedisCfg struct {
	Host string  `toml:"host"`
	Port int     `toml:"port"`
	Db   int     `toml:"db"`
	Pool PoolCfg `toml:"pool"`
}

func (p *RedisCfg) Server() string {
	return fmt.Sprintf("%s:%d", p.Host, p.Port)
}

type PoolCfg struct {
	MaxIdle int `toml:"max_idle"`
	MaxOpen int `toml:"max_open"`
}

type ElasticsearchCfg struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

func (p *Configuration) Store(file string) error {
	fi, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fi.Close()

	end := toml.NewEncoder(fi)
	p.SecretsS = ToBase64(p.Secrets)
	return end.Encode(p)

}

func (p *Configuration) Load(file string) error {
	var err error
	if _, err = toml.DecodeFile(file, p); err != nil {
		return err
	}
	p.Secrets, err = FromBase64(p.SecretsS)
	return err
}

func (p *Configuration) IsProduction() bool {
	return p.Env == "production"
}

func (p *Configuration) OpenRedis() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     p.Redis.Pool.MaxIdle,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%d", p.Redis.Host, p.Redis.Port))
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func (p *Configuration) OpenDb() (*gorm.DB, error) {
	db, err := gorm.Open(p.Database.Dialect, p.Database.Url)
	if err != nil {
		return nil, err
	}
	db.LogMode(!p.IsProduction())

	db.DB().SetMaxIdleConns(p.Database.Pool.MaxIdle)
	db.DB().SetMaxOpenConns(p.Database.Pool.MaxOpen)
	if err = db.DB().Ping(); err != nil {
		return nil, err
	}
	return &db, nil
}

//==============================================================================
func Load(c *cli.Context) (*Configuration, error) {

	var cfg Configuration
	env := c.String("environment")
	if err := cfg.Load(fmt.Sprintf("config/%s/settings.toml", env)); err == nil {
		cfg.Env = env
		return &cfg, nil
	} else {
		return nil, err
	}
}
