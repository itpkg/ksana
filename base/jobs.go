package base

import (
	"github.com/jrallison/go-workers"
)

func (p *BaseEngine) Deploy() {
}

func (p *BaseEngine) Cron() map[string]func() {
	return map[string]func(){
		"0 0 3 * * *": func() {
			//todo
		},
	}
}

func (p *BaseEngine) Worker() {

	workers.Process("email",
		func(message *workers.Msg) {

		},
		p.Cfg.Workers["email"])
}
