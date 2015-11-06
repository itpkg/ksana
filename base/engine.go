package base

import (
	"github.com/itpkg/ksana/i18n"
	"github.com/itpkg/ksana/ioc"
	"github.com/itpkg/ksana/mux"
	"github.com/itpkg/ksana/settings"
)

type SiteEngine struct {
	Settings settings.Store `inject:""`
	I18n     i18n.Store     `inject:""`
}

type UsersEngine struct {
	Dao *Dao `inject:""`
}

func init() {
	se := SiteEngine{}
	ue := UsersEngine{}
	ioc.Use(&se, &ue)

	mux.Register(&se, &ue)
}
