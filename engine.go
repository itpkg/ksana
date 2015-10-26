package ksana

type Engine interface {
	Router()
	Migrate()
	Job()
	Deploy()
	Shell()
}
