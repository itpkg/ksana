package ksana

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Configuration struct {
	Env           string           `toml:"-"`
	Http          HttpCfg          `toml:"http"`
	Database      DatabaseCfg      `toml:"database"`
	Redis         RedisCfg         `toml:"redis"`
	Elasticsearch ElasticsearchCfg `toml:"elasticsearch"`
}

type HttpCfg struct {
	Host    string `toml:"host"`
	Port    int    `toml:"port"`
	Secrets string `toml:"secrets"`
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
	return end.Encode(p)
}

func (p *Configuration) Load(file string) error {
	_, err := toml.DecodeFile(file, p)
	return err
}
