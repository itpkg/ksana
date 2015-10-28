package orm

import (
	"time"
)

type Dialect interface {
	CreateDatabase(name string) string
	DropDatabase(name string) string
	CreateTable(name string, columns ...string) string
	DropTable(name string) string
	AddIndex(table string, columns ...string) string
	AddUniqueIndex(table string, columns ...string) string
	DropIndex(table string, columns ...string) string

	Id() string
	CreatedAt() string
	UpdatedAt() string

	String(name string, length uint, nullable bool, def_val interface{}) string
	Chars(name string, length uint, nullable bool, def_val interface{}) string
	Text(name string, nullable bool) string
	Time(name string, nullable bool, def_val *time.Time) string
	Date(name string, nullable bool, def_val *time.Time) string
	Timestamp(name string, nullable bool, def_val *time.Time) string
	Bool(name string, nullable bool, def_val interface{}) string
	Int(name string, nullable bool, def_val interface{}) string
	Long(name string, nullable bool, def_val interface{}) string
	Bytes(name string, nullable bool) string
}
