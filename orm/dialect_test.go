package orm_test

import (
	"testing"
	"time"

	ko "github.com/itpkg/ksana/orm"
)

func TestPostgresqlDialect(t *testing.T) {
	run_dialect(t, &ko.PostgresqlDialect{})
}

func run_dialect(t *testing.T, d ko.Dialect) {
	now := time.Now()
	for _, s := range []string{
		d.CreateDatabase("ksana_test"),
		d.DropDatabase("ksana_test"),
		d.CreateTable("t1", d.Id(), d.CreatedAt(), d.UpdatedAt(), d.Int("c1", false, 123), d.String("c2", 255, false, "==="), d.Timestamp("c3", false, &now)),
		d.AddIndex("t1", "c1", "c2"),
		d.AddUniqueIndex("t1", "c1", "c2", "c3"),
		d.DropIndex("t1", "c1", "c2"),
		d.DropTable("t1"),
	} {
		t.Logf(s)
	}
}
