package mux

import (
	"net/http"
	"regexp"
)

type Route struct {
	Method   string
	Pattern  *regexp.Regexp
	Handlers []HttpHandler
	names    []string
}

func (p *Route) Match(req *http.Request) map[string]string {

	match := p.Pattern.FindStringSubmatch(req.URL.Path)
	if match == nil {
		return nil
	}
	params := make(map[string]string)
	for i, n := range p.Pattern.SubexpNames() {
		if i > 0 {
			params[n] = match[i]
		}
	}
	return params
}
