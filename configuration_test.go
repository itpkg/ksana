package ksana_test

import (
	"testing"

	ks "github.com/itpkg/ksana"
)

const cfg_file = "config.toml"
const port = 8080

func TestCfgStore(t *testing.T) {
	buf, _ := ks.RandomBytes(512)
	cfg := ks.Configuration{
		Http: ks.HttpCfg{
			Host:    "localhost",
			Port:    port,
			Secrets: ks.ToBase64(buf),
		},
		Database: ks.DatabaseCfg{
			Dialect: "postgres",
			Url:     "user=postgres dbname=itpkg sslmode=disable",
			Pool: ks.PoolCfg{
				MaxIdle: 6,
				MaxOpen: 180,
			},
		},
		Redis: ks.RedisCfg{
			Host: "localhost",
			Port: 6379,
			Db:   0,
			Pool: ks.PoolCfg{
				MaxIdle: 4,
				MaxOpen: 120,
			},
		},
		Elasticsearch: ks.ElasticsearchCfg{
			Host: "localhost",
			Port: 9200,
		},
		Workers: map[string]int{"aaa": 1, "bbb": 2, "ccc": 3},
	}

	if err := cfg.Store(cfg_file); err != nil {
		t.Errorf("store cfg error: %v", err)
	}

}

func TestCfgLoad(t *testing.T) {
	var cfg ks.Configuration
	if err := cfg.Load(cfg_file); err != nil {
		t.Errorf("load cfg error: %v", err)
	}
	t.Logf("cfg redis: %v", cfg.Redis)
	if cfg.Http.Port != port {
		t.Errorf("bad value: %d vs %d", cfg.Http.Port, port)
	}
}
