package mux

import (
	"net/http"
)

type Router struct {
	Routes   []*Route
	Handlers []HttpHandler
}

func (p *Router) Get(pattern string, handlers ...HttpHandler) {
	p.Add("GET", pattern, handlers...)
}

func (p *Router) Post(pattern string, handlers ...HttpHandler) {
	p.Add("POST", pattern, handlers...)
}

func (p *Router) Put(pattern string, handlers ...HttpHandler) {
	p.Add("PUT", pattern, handlers...)
}

func (p *Router) Patch(pattern string, handlers ...HttpHandler) {
	p.Add("PATCH", pattern, handlers...)
}

func (p *Router) Delete(pattern string, handlers ...HttpHandler) {
	p.Add("DELETE", pattern, handlers...)
}

func (p *Router) Add(method, pattern string, handlers ...HttpHandler) {
	p.Routes = append(
		p.Routes,
		&Route{
			Method:   method,
			Pattern:  pattern,
			Handlers: handlers,
		},
	)
}

func (p *Router) ServeHTTP(http.ResponseWriter, *http.Request) {
	//todo
}
