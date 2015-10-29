package orm

import (
	"database/sql"
)

type Handler func(vls ...interface{}) error

type Db struct {
	db *sql.DB
}

func (p *Db) Open(driver, source string) error {
	var err error
	if p.db, err = sql.Open(driver, source); err != nil {
		return err
	}
	err = p.db.Ping()
	return err
}

func (p *Db) SetMaxOpen(n int) {
	p.db.SetMaxOpenConns(n)
}

func (p *Db) SetMaxIdle(n int) {
	p.db.SetMaxIdleConns(n)
}

func (p *Db) List(query string, fn Handler, args ...interface{}) ([]interface{}, error) {
	rows, err := p.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {

	}
	return nil, nil

}
func (p *Db) Get() {
}
func (p *Db) Del() {
}
func (p *Db) Set() {
}

func (p *Db) Exec(query string, args ...interface{}) error {
	_, err := p.db.Exec(query, args...)
	return err
}
