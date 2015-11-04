package settings

import (
	orm "github.com/itpkg/ksana/orm"
	utils "github.com/itpkg/ksana/utils"
)

type DatabaseStore struct {
	db  *orm.Db
	aes *utils.Aes
}

func (p *DatabaseStore) Set(key string, val interface{}, encrypt bool) error {
	buf, err := utils.ToJson(val)
	if err != nil {
		return err
	}
	if encrypt {
		buf, err = p.aes.Encrypt(buf)
		if err != nil {
			return err
		}
	}

	row := p.db.Get("settings.count", key)

	var c int

	if err = row.Scan(&c); err != nil {
		return err
	}
	if c > 0 {
		_, err = p.db.Exec("settings.update", encrypt, buf, key)
	} else {
		_, err = p.db.Exec("settings.add", key, encrypt, buf)

	}
	return err
}

func (p *DatabaseStore) Get(key string, val interface{}) error {
	row := p.db.Get("settings.get", key)

	var buf []byte
	var enc bool

	err := row.Scan(&enc, &buf)
	if err != nil {
		return err
	}
	if enc {
		buf, err = p.aes.Decrypt(buf)
		if err != nil {
			return err
		}
	}
	return utils.FromJson(buf, val)

}

//==============================================================================

func NewDatabaseStore(db *orm.Db, aes *utils.Aes) (Store, error) {
	ds := DatabaseStore{db: db, aes: aes}
	err := db.Load(utils.PkgRoot(&ds))
	if err == nil {
		return &ds, nil
	} else {
		return nil, err
	}
}
