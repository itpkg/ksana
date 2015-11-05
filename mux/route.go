package mux

import (
	"net/http"
	"regexp"
)

type Route struct {
	method      string
	pattern     *regexp.Regexp
	handlers    []Handler
	middlewares []Middleware
}

func (p *Route) Match(req *http.Request) bool {
	return p.pattern.MatchString(req.URL.Path)
}

func (p *Route) Parse(req *http.Request) map[string]interface{} {

	match := p.pattern.FindStringSubmatch(req.URL.Path)
	if match == nil {
		return nil
	}
	params := make(map[string]interface{})
	for i, n := range p.pattern.SubexpNames() {
		if i > 0 {
			params[n] = match[i]
		}
	}
	return params
}
