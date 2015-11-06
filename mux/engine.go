package mux

import (
	"github.com/itpkg/ksana/atom"
	"github.com/itpkg/ksana/sitemap"
)

type Engine interface {
	Mount(*Router)
	Sitemap() []*sitemap.Item
	Atom() []*atom.Entry
}

var engines = make([]Engine, 0)

func Register(ens ...Engine) {
	engines = append(engines, ens...)
}

func Loop(fn func(Engine) error) error {
	for _, en := range engines {
		if err := fn(en); err != nil {
			return err
		}
	}
	return nil
}
