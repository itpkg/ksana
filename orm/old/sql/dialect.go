package sql

import(
	"time"
)

type Dialect interface{
	CreateDatabase(name string) string
	DropDatabase(name string) string
	CreateTable(name string, columns ...string) string
	DropTable(name string) string
	CreateIndex(table string, unique bool, columns ...string) string
	DropIndex(table string, columns ...string) string
	//-------------------------
	Id() string
	Created() string
	Updated() string
	Deleted() string
	
	String(name string, length int, nullable bool, def_val string) string
	FixString(name string, length int, nullable bool, def_val string) string
	BigString(name string, nullable bool) string

	Bytes(name string, length int, nullable bool)
	BigBytes(name string, nullable bool)

	Float(name string, nullable bool, def_val float32)
	Double(name string, nullable bool, def_val float64)

	Int(name string, nullable bool, def_val int)
	Long(name string, nullable bool, def_val int)

	Bool(name string, nullable bool, def_val bool)

	Date(name string, nullable_bool, def_val *time.Time)
	Time(name string, nullable_bool, def_val *time.Time)
	Datetime(name string, nullable_bool, def_val *time.Time)
}
