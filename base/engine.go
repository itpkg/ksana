package base

import (
	ks "github.com/itpkg/ksana"
	"github.com/jinzhu/gorm"
)

type BaseEngine struct {
	Db  *gorm.DB          `inject:""`
	Dao *Dao              `inject:""`
	Cfg *ks.Configuration `inject:""`
}

//==============================================================================
func init() {
	var en ks.Engine
	en = &BaseEngine{}
	ks.Use(en)
}
