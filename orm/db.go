package orm

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/itpkg/ksana/logging"
	"github.com/itpkg/ksana/utils"
)

type Db struct {
	db         *sql.DB
	cfg        *Configuration
	mapper     map[string]string
	migrations []*Migration
	Logger     logging.Logger `inject:""`
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

func (p *Db) query(name string) string {
	q := p.mapper[name]
	p.Logger.Debug(q)
	return q

}

func (p *Db) Exec(query string, args ...interface{}) (sql.Result, error) {
	return p.db.Exec(p.query(query), args...)
}

func (p *Db) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return p.db.Query(p.query(query), args...)
}

func (p *Db) Get(query string, args ...interface{}) *sql.Row {
	return p.db.QueryRow(p.query(query), args...)
}

func (p *Db) ExecT(tx *sql.Tx, query string, args ...interface{}) (sql.Result, error) {
	return tx.Exec(p.query(query), args...)
}

func (p *Db) GetT(tx *sql.Tx, query string, args ...interface{}) *sql.Row {
	return tx.QueryRow(p.query(query), args...)
}

func (p *Db) Migrate() error {
	return p.Commit(func(tx *sql.Tx) (err error) {
		if _, err = p.ExecT(tx, "schema_migrations.check"); err != nil {
			return
		}
		for _, m := range p.migrations {
			var c int
			row := p.GetT(tx, "schema_migrations.count", m.Id)
			if err = row.Scan(&c); err != nil {
				return
			}
			if c > 0 {
				continue
			}
			for _, s := range m.Up {
				p.Logger.Debug(s)
				if _, err = tx.Exec(s); err != nil {
					return
				}
			}
			if _, err = p.ExecT(tx, "schema_migrations.add", m.Id); err != nil {
				return
			}
		}

		return nil
	})

}

func (p *Db) Rollback() error {
	return p.Commit(func(tx *sql.Tx) (err error) {
		if _, err = p.ExecT(tx, "schema_migrations.check"); err != nil {
			return
		}
		var id int
		var ver string
		row := p.GetT(tx, "schema_migrations.last")
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
					p.Logger.Debug(s)
					if _, err = tx.Exec(s); err != nil {
						return
					}
				}
			}
		}
		_, err = p.ExecT(tx, "schema_migrations.remove", id)

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

func (p *Db) load(dir string) error {
	err := p.list_files(fmt.Sprintf("%s/mappers/%s", dir, p.cfg.Driver), func(d, n string) error {
		tmp := make(map[string]string, 0)
		if e := utils.FromToml(fmt.Sprintf("%s/%s", d, n), &tmp); e != nil {
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
		if e := utils.FromToml(fmt.Sprintf("%s/%s", d, n), &mig); e != nil {
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

func init() {
	Register(utils.PkgRoot((*Db)(nil)))
}

func Open(dir string) (*Db, error) {
	cfg := Configuration{}
	if err := utils.FromToml(fmt.Sprintf("%s/database.toml", dir), &cfg); err != nil {
		return nil, err
	}

	//--------
	db, err := sql.Open(cfg.Driver, cfg.Source())

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

	for _, m := range modules {
		if err = rdb.load(m); err != nil {
			return nil, err
		}
	}

	return &rdb, nil

}
