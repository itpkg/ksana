package mux

type Engine interface {
	Mount(*Router)
}

var engines = make([]Engine, 0)

func Register(ens ...Engine) {
	engines = append(engines, ens...)
}
