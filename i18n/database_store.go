package i18n

import (
	"fmt"
	"os"
	"reflect"

	orm "github.com/itpkg/ksana/orm"
)

type DatabaseStore struct {
	db *orm.Db
}

func (p *DatabaseStore) Set(lang, code, msg string) error {
	row := p.db.Get("i18n.count", lang, code)

	var c int
	var err error

	if err = row.Scan(&c); err != nil {
		return err
	}
	if c > 0 {
		_, err = p.db.Exec("i18n.update", msg, lang, code)
	} else {
		_, err = p.db.Exec("i18n.add", lang, code, msg)

	}
	return err
}

func (p *DatabaseStore) Get(lang, code string) (string, error) {
	row := p.db.Get("i18n.get", lang, code)

	var msg string

	err := row.Scan(&msg)
	return msg, err
}

func (p *DatabaseStore) Loop(lang string, fn func(_, _ string) error) error {
	rows, err := p.db.Query("i18n.all", lang)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var code string
		var msg string
		if err = rows.Scan(&code, &msg); err != nil {
			return err
		}
		if err = fn(code, msg); err != nil {
			return err
		}
	}
	return nil
}

//==============================================================================

func NewDatabaseStore(db *orm.Db) (Store, error) {
	ds := DatabaseStore{db: db}
	err := db.Load(fmt.Sprintf("%s/src/%s", os.Getenv("GOPATH"), reflect.TypeOf(&ds).Elem().PkgPath()))
	if err == nil {
		return &ds, nil
	} else {
		return nil, err
	}
}
