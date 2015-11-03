package orm

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/BurntSushi/toml"
)

type Configuration struct {
	Driver   string                 `toml:"driver"`
	Host     string                 `toml:"host"`
	Port     int                    `toml:"port"`
	User     string                 `toml:"user"`
	Password string                 `toml:"password"`
	Name     string                 `toml:"name"`
	MaxOpen  int                    `toml:"max_open"`
	MaxIdle  int                    `toml:"max_idle"`
	Extra    map[string]interface{} `toml:"extra"`
}

func (p *Configuration) Source() string {
	ex := make([]string, 0)
	for k, v := range p.Extra {
		ex = append(ex, fmt.Sprintf("%s=%v", k, v))
	}
	return fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s?%s",
		p.Driver,
		p.User,
		p.Password,
		p.Host,
		p.Port,
		p.Name,
		strings.Join(ex, "&"),
	)
}

//==============================================================================
type Db struct {
	db         *sql.DB
	mapper     map[string]string
	migrations []*Migration
}

func (p *Db) Status(w io.Writer) {
	fmt.Fprintf(w, "=== MAPPER ===\n")
	for k, v := range p.mapper {
		fmt.Fprintf(w, "%s: %s\n", k, v)
	}

	for i, m := range p.migrations {
		fmt.Fprintf(w, "=== %d: %s ===\n", i, m.Id)
		fmt.Fprintf(w, "[UP]\n%s\n", strings.Join(m.Up, "\n"))
		fmt.Fprintf(w, "[DOWN]\n%s\n", strings.Join(m.Down, "\n"))
	}
}

func (p *Db) Open(dir string) error {
	cfg := Configuration{}
	_, err := toml.DecodeFile(fmt.Sprintf("%s/database.toml", dir), &cfg)
	if err != nil {
		return err
	}

	//--------
	p.db, err = sql.Open(cfg.Driver, cfg.Source())

	if err == nil {
		p.db.SetMaxOpenConns(cfg.MaxOpen)
		p.db.SetMaxIdleConns(cfg.MaxIdle)
		err = p.db.Ping()
	}
	if err != nil {
		return err
	}

	//--------
	p.mapper = make(map[string]string, 0)
	err = p.list_files(fmt.Sprintf("%s/mapper/%s", dir, cfg.Driver), func(d, n string) error {
		tmp := make(map[string]string, 0)
		if _, e := toml.DecodeFile(fmt.Sprintf("%s/%s", d, n), &tmp); e != nil {
			return e
		}
		for k, v := range tmp {
			p.mapper[k] = v
		}
		return nil
	})
	if err != nil {
		return err
	}

	//---------
	p.migrations = make([]*Migration, 0)
	err = p.list_files(fmt.Sprintf("%s/migrations/%s", dir, cfg.Driver), func(d, n string) error {
		mig := Migration{}
		if _, e := toml.DecodeFile(fmt.Sprintf("%s/%s", d, n), &mig); e != nil {
			return e
		}
		mig.Id = n[0 : len(n)-5]
		p.migrations = append(p.migrations, &mig)
		return nil
	})

	return err
}

func (p *Db) list_files(dir string, fn func(_, _ string) error) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, f := range files {
		if err = fn(dir, f.Name()); err != nil {
			return err
		}
	}
	return nil

}
