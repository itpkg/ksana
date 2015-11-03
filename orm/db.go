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
	cfg        *Configuration
	mapper     map[string]string
	migrations []*Migration
}

func (p *Db) Commit(fn func(_ *sql.Tx) error) (err error) {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			err = tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	return fn(tx)
}

func (p *Db) Exec(tx *sql.Tx, query string, args ...interface{}) (sql.Result, error) {
	//todo logger
	return tx.Exec(p.mapper[query], args...)
}

func (p *Db) Select(tx *sql.Tx, query string, args ...interface{}) *sql.Row {
	//todo logger
	return tx.QueryRow(p.mapper[query], args...)
}

func (p *Db) Migrate() error {
	return p.Commit(func(tx *sql.Tx) (err error) {
		if _, err = p.Exec(tx, "schema_migrations.check"); err != nil {
			return
		}
		for _, m := range p.migrations {
			var c int
			row := p.Select(tx, "schema_migrations.count", m.Id)
			if err = row.Scan(&c); err != nil {
				return
			}
			if c > 0 {
				continue
			}
			for _, s := range m.Up {
				if _, err = tx.Exec(s); err != nil {
					return
				}
			}
			if _, err = p.Exec(tx, "schema_migrations.add", m.Id); err != nil {
				return
			}
		}

		return nil
	})

}

func (p *Db) Rollback() error {
	return p.Commit(func(tx *sql.Tx) (err error) {
		if _, err = p.Exec(tx, "schema_migrations.check"); err != nil {
			return
		}
		var id int
		var ver string
		row := p.Select(tx, "schema_migrations.last")
		err = row.Scan(&id, &ver)
		if err == sql.ErrNoRows {
			return nil
		}
		if err != nil {
			return
		}
		for _, m := range p.migrations {
			if ver == m.Id {
				for _, s := range m.Down {
					if _, err = tx.Exec(s); err != nil {
						return
					}
				}
			}
		}
		_, err = p.Exec(tx, "schema_migrations.remove", id)

		return
	})
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

func (p *Db) Load(dir string) error {
	err := p.list_files(fmt.Sprintf("%s/mapper/%s", dir, p.cfg.Driver), func(d, n string) error {
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
	err = p.list_files(fmt.Sprintf("%s/migrations/%s", dir, p.cfg.Driver), func(d, n string) error {
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

//==============================================================================

func Open(dir string) (*Db, error) {
	cfg := Configuration{}
	_, err := toml.DecodeFile(fmt.Sprintf("%s/database.toml", dir), &cfg)
	if err != nil {
		return nil, err
	}

	//--------
	var db *sql.DB
	db, err = sql.Open(cfg.Driver, cfg.Source())

	if err == nil {
		db.SetMaxOpenConns(cfg.MaxOpen)
		db.SetMaxIdleConns(cfg.MaxIdle)
		err = db.Ping()
	}
	if err != nil {
		return nil, err
	}

	rdb := Db{
		db:         db,
		cfg:        &cfg,
		mapper:     make(map[string]string, 0),
		migrations: make([]*Migration, 0),
	}

	err = rdb.Load(dir)
	if err == nil {
		switch cfg.Driver {
		case "postgres":
			rdb.mapper["schema_migrations.check"] = "CREATE TABLE IF NOT EXISTS schema_migrations(id SERIAL NOT NULL PRIMARY KEY, version VARCHAR(255) NOT NULL UNIQUE, created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)"
			rdb.mapper["schema_migrations.last"] = "SELECT id, version FROM schema_migrations ORDER BY id DESC LIMIT 1"
			rdb.mapper["schema_migrations.remove"] = "DELETE FROM schema_migrations WHERE id = $1"
			rdb.mapper["schema_migrations.count"] = "SELECT count(*) FROM schema_migrations WHERE version = $1"
			rdb.mapper["schema_migrations.add"] = "INSERT INTO schema_migrations(version) VALUES($1)"
		case "mysql":
			//todo
		}
		return &rdb, nil
	} else {
		return nil, err
	}
}
