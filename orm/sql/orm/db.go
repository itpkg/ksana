package orm

import (
	"database/sql"
)

type Db struct {
	driver string
	db     *sql.DB
}

func (p *Db) From(table string) *Query {
	return &Query{table: table}
}
