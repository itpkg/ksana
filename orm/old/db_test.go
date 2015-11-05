package orm_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/BurntSushi/toml"
	kl "github.com/itpkg/ksana/logging"
	ko "github.com/itpkg/ksana/orm"
	_ "github.com/lib/pq"
)

var cfg_f = "config/test/database.toml"
var cfg = ko.Configuration{
	Driver:   "postgres",
	Host:     "localhost",
	Port:     5432,
	User:     "postgres",
	Password: "",
	Name:     "ksana_test",
	Extra:    map[string]interface{}{"sslmode": "disable"},
	MaxOpen:  120,
	MaxIdle:  6,
}

var mapper1_f = "tmp/mappers/postgres/1.toml"
var mapper1 = map[string]string{
	"current.time": "SELECT CURRENT_TIME",
	"current.date": "SELECT CURRENT_DATE",
}
var mapper2_f = "tmp/mappers/postgres/2.toml"
var mapper2 = map[string]string{
	"current.timestamp": "SELECT CURRENT_TIMESTAMP",
}

var mig1 = ko.Migration{
	Id: "111-create_1",
	Up: []string{
		"create table t11(id int)",
		"create table t12(id int)",
	},
	Down: []string{
		"drop table t11",
		"drop table t12",
	},
}

var mig2 = ko.Migration{
	Id: "222-create_2",
	Up: []string{
		"create table t21(id int, \"key\" varchar(255))",
		"create table t22(id int)",
	},
	Down: []string{
		"drop table t21",
		"drop table t22",
	},
}

func write(t *testing.T, f string, v interface{}) {
	fi, err := os.Create(f)
	defer fi.Close()

	if err == nil {
		end := toml.NewEncoder(fi)
		err = end.Encode(v)
	}

	if err != nil {
		t.Errorf("error on write: %v", err)
	}
}

func TestConfig(t *testing.T) {
	write(t, cfg_f, &cfg)
	write(t, mapper1_f, mapper1)
	write(t, mapper2_f, mapper2)

	write(t, fmt.Sprintf("tmp/migrations/postgres/%s.toml", mig1.Id), mig1)
	write(t, fmt.Sprintf("tmp/migrations/postgres/%s.toml", mig2.Id), mig2)
}

func TestOpen(t *testing.T) {

	ko.Register("tmp")
	db, err := ko.Open("test")
	if err != nil {
		t.Errorf("error on open: %v", err)
	}
	db.Logger = kl.NewStdoutLogger(kl.DEBUG)

	db.Status(os.Stdout)

	if err = db.Migrate(); err != nil {
		t.Errorf("error on migrate: %v", err)
	}

	if err = db.Rollback(); err != nil {
		t.Errorf("error on rollback: %v", err)
	}

}
