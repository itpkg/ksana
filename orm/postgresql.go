package orm

import (
	"fmt"
	"strings"
	"time"
)

type PostgresqlDialect struct {
}

func (p *PostgresqlDialect) CreateDatabase(name string) string {
	return fmt.Sprintf("CREATE DATABASE %s ENCODING='UTF8'", name)
}
func (p *PostgresqlDialect) DropDatabase(name string) string {
	return fmt.Sprintf("DROP DATABASE IF EXISTS %s", name)
}

func (p *PostgresqlDialect) CreateTable(name string, columns ...string) string {
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(%s)", name, strings.Join(columns, ", "))
}

func (p *PostgresqlDialect) DropTable(name string) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %s", name)
}

func (p *PostgresqlDialect) AddIndex(table string, columns ...string) string {
	return fmt.Sprintf("CREATE INDEX %s ON %s (%s)", p.index_name(table, columns...), table, strings.Join(columns, ", "))
}

func (p *PostgresqlDialect) AddUniqueIndex(table string, columns ...string) string {
	return fmt.Sprintf("CREATE UNIQUE INDEX %s ON %s (%s)", p.index_name(table, columns...), table, strings.Join(columns, ", "))
}

func (p *PostgresqlDialect) DropIndex(table string, columns ...string) string {
	return fmt.Sprintf("DROP INDEX IF EXISTS %s", p.index_name(table, columns...))
}

func (p *PostgresqlDialect) index_name(table string, columns ...string) string {
	return fmt.Sprintf("idx_%s_%s", table, strings.Join(columns, "_"))
}

func (p *PostgresqlDialect) Id() string {
	return fmt.Sprintf("ID SERIAL PRIMARY KEY")
}

func (p *PostgresqlDialect) String(name string, length uint, nullable bool, def_val interface{}) string {
	ns := ""
	if !nullable {
		ns = " NOT NULL"
	}
	dv := ""
	if def_val != nil {
		dv = fmt.Sprintf(" DEFAULT '%s'", def_val.(string))
	}
	return fmt.Sprintf("%s VARCHAR(%d)%s%s", name, length, ns, dv)
}

func (p *PostgresqlDialect) Chars(name string, length uint, nullable bool, def_val interface{}) string {
	ns := ""
	if !nullable {
		ns = " NOT NULL"
	}
	dv := ""
	if def_val != nil {
		dv = fmt.Sprintf(" DEFAULT '%s'", def_val.(string))
	}
	return fmt.Sprintf("%s CHAR(%d)%s%s", name, length, ns, dv)
}

func (p *PostgresqlDialect) CreatedAt() string {
	return "created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP"
}

func (p *PostgresqlDialect) UpdatedAt() string {
	return "updated_at TIMESTAMP NOT NULL"
}

func (p *PostgresqlDialect) Text(name string, nullable bool) string {
	ns := ""
	if !nullable {
		ns = " NOT NULL"
	}
	return fmt.Sprintf("%s TEXT%s", name, ns)
}

func (p *PostgresqlDialect) Time(name string, nullable bool, def_val *time.Time) string {
	ns := ""
	if !nullable {
		ns = " NOT NULL"
	}
	dv := ""
	if def_val != nil {
		dv = fmt.Sprintf(" DEFAULT '%s'", def_val.Format("15:04:05"))
	}
	return fmt.Sprintf("%s TIME%s%s", name, ns, dv)
}

func (p *PostgresqlDialect) Date(name string, nullable bool, def_val *time.Time) string {
	ns := ""
	if !nullable {
		ns = " NOT NULL"
	}
	dv := ""
	if def_val != nil {
		dv = fmt.Sprintf(" DEFAULT '%s'", def_val.Format("2006-01-02"))
	}
	return fmt.Sprintf("%s DATE%s%s", name, ns, dv)
}

func (p *PostgresqlDialect) Timestamp(name string, nullable bool, def_val *time.Time) string {
	ns := ""
	if !nullable {
		ns = " NOT NULL"
	}
	dv := ""
	if def_val != nil {
		dv = fmt.Sprintf(" DEFAULT '%s'", def_val.Format("2006-01-02 15:04:05"))
	}
	return fmt.Sprintf("%s TIMESTAMP%s%s", name, ns, dv)
}

func (p *PostgresqlDialect) Bool(name string, nullable bool, def_val interface{}) string {
	ns := ""
	if !nullable {
		ns = " NOT NULL"
	}
	dv := ""
	if def_val != nil {
		dv = fmt.Sprintf(" DEFAULT '%v'", def_val.(bool))
	}
	return fmt.Sprintf("%s BOOLEAN%s%s", name, ns, dv)
}

func (p *PostgresqlDialect) Int(name string, nullable bool, def_val interface{}) string {
	ns := ""
	if !nullable {
		ns = " NOT NULL"
	}
	dv := ""
	if def_val != nil {
		dv = fmt.Sprintf(" DEFAULT %d", def_val.(int))
	}
	return fmt.Sprintf("%s INT%s%s", name, ns, dv)
}

func (p *PostgresqlDialect) Long(name string, nullable bool, def_val interface{}) string {
	ns := ""
	if !nullable {
		ns = " NOT NULL"
	}
	dv := ""
	if def_val != nil {
		dv = fmt.Sprintf(" DEFAULT '%v'", def_val.(int))
	}
	return fmt.Sprintf("%s BIGINT%s%s", name, ns, dv)
}

func (p *PostgresqlDialect) Bytes(name string, nullable bool) string {
	ns := ""
	if !nullable {
		ns = " NOT NULL"
	}
	return fmt.Sprintf("%s BYTEA%s", name, ns)
}
