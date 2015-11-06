package base

import (
	"time"

	"github.com/itpkg/ksana/orm"
	"github.com/itpkg/ksana/utils"
)

type User struct {
	Id           uint
	Username     string
	Email        string
	Token        string
	Password     []byte
	ProviderId   string
	ProviderType string
	Status       string
	Created      time.Time
	Updated      time.Time
}

type Log struct {
	Id      uint
	UserId  uint
	User    *User
	Message string
	Type    string
	Created time.Time
}

//==============================================================================

type Dao struct {
	Db *orm.Db `inject:""`
}

//==============================================================================
func init() {
	orm.Register(utils.PkgRoot((*Dao)(nil)))
}
