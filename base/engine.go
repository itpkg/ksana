package base

import (
	"github.com/itpkg/ksana/ioc"
	"github.com/itpkg/ksana/mux"
)

type SiteEngine struct {
}

type UsersEngine struct {
}

func init() {
	se := SiteEngine{}
	ue := UsersEngine{}
	ioc.Use(&se, &ue)

	mux.Register(&se, &ue)
}
