package i18n

import (
	"github.com/itpkg/ksana/orm"
	"github.com/itpkg/ksana/utils"
)

type DatabaseStore struct {
	Db *orm.Db `inject:""`
}

func (p *DatabaseStore) Set(lang, code, msg string) error {
	row := p.Db.Get("i18n.count", lang, code)

	var c int
	var err error

	if err = row.Scan(&c); err != nil {
		return err
	}
	if c > 0 {
		_, err = p.Db.Exec("i18n.update", msg, lang, code)
	} else {
		_, err = p.Db.Exec("i18n.add", lang, code, msg)

	}
	return err
}

func (p *DatabaseStore) Get(lang, code string) (string, error) {
	row := p.Db.Get("i18n.get", lang, code)

	var msg string

	err := row.Scan(&msg)
	return msg, err
}

func (p *DatabaseStore) Loop(lang string, fn func(_, _ string) error) error {
	rows, err := p.Db.Query("i18n.all", lang)
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
func init() {
	orm.Register(utils.PkgRoot((*DatabaseStore)(nil)))
}
