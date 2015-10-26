package ksana

import (
	"io"

	"github.com/BurntSushi/toml"
)

type Configuration struct {
	Http          HttpCfg          `toml:"http"`
	Database      DatabaseCfg      `toml:"database"`
	Redis         RedisCfg         `toml:"redis"`
	Elasticsearch ElasticsearchCfg `toml:"elasticsearch"`
}

type HttpCfg struct {
	Port    int
	Secrets string
}

type DatabaseCfg struct {
	Dialect string
	Url     string
	Pool    PoolCfg `toml:"pool"`
}

type RedisCfg struct {
	Host string
	Port int
	Db   int
	Pool PoolCfg `toml:"pool"`
}

type PoolCfg struct {
	MaxIdle int
	MaxOpen int
}

type ElasticsearchCfg struct {
	Host string
	Port int
}

func (p *Configuration) Store(w io.Writer) error {
	end := toml.NewEncoder(w)
	return end.Encode(p)
}

func (p *Configuration) Load(data string) error {
	_, err := toml.Decode(data, p)
	return err
}
