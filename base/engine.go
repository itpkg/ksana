package base

import (
	"github.com/itpkg/ksana/atom"
	"github.com/itpkg/ksana/i18n"
	"github.com/itpkg/ksana/ioc"
	"github.com/itpkg/ksana/mux"
	"github.com/itpkg/ksana/settings"
	"github.com/itpkg/ksana/sitemap"
)

type SiteEngine struct {
	Settings settings.Store `inject:""`
	I18n     i18n.Store     `inject:""`
}

func (p *SiteEngine) Sitemap() []*sitemap.Item {
	return nil
}
func (p *SiteEngine) Atom() []*atom.Entry {
	return nil
}

type UsersEngine struct {
	Dao *Dao `inject:""`
}

func (p *UsersEngine) Sitemap() []*sitemap.Item {
	return nil
}
func (p *UsersEngine) Atom() []*atom.Entry {
	return nil
}

func init() {
	se := SiteEngine{}
	ue := UsersEngine{}
	ioc.Use(&se, &ue)

	mux.Register(&se, &ue)
}
