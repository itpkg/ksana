package settings

import (
	"github.com/itpkg/ksana/orm"
	"github.com/itpkg/ksana/utils"
)

type DatabaseStore struct {
	Db  *orm.Db    `inject:""`
	Aes *utils.Aes `inject:""`
}

func (p *DatabaseStore) Set(key string, val interface{}, encrypt bool) error {
	buf, err := utils.ToJson(val)
	if err != nil {
		return err
	}
	if encrypt {
		buf, err = p.Aes.Encrypt(buf)
		if err != nil {
			return err
		}
	}

	row := p.Db.Get("settings.count", key)

	var c int

	if err = row.Scan(&c); err != nil {
		return err
	}
	if c > 0 {
		_, err = p.Db.Exec("settings.update", encrypt, buf, key)
	} else {
		_, err = p.Db.Exec("settings.add", key, encrypt, buf)

	}
	return err
}

func (p *DatabaseStore) Get(key string, val interface{}) error {
	row := p.Db.Get("settings.get", key)

	var buf []byte
	var enc bool

	err := row.Scan(&enc, &buf)
	if err != nil {
		return err
	}
	if enc {
		buf, err = p.Aes.Decrypt(buf)
		if err != nil {
			return err
		}
	}
	return utils.FromJson(buf, val)

}

//==============================================================================
func init() {
	orm.Register(utils.PkgRoot((*DatabaseStore)(nil)))
}
