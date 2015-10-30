package mux

import (
	"net/http"
	"regexp"
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
			Pattern:  regexp.MustCompile(pattern),
			Handlers: handlers,
		},
	)
}

func (p *Router) ServeHTTP(wrt http.ResponseWriter, req *http.Request) {
	for _, rt := range p.Routes {
		if rt.Match(req) != nil {
			ctx := Context{
				Request: req,
				Writer:  wrt,
				params:  make(map[string]interface{}, 0),
			}

			for _, h := range append(p.Handlers, rt.Handlers...) {
				if c, e := h(&ctx); e != nil {
					wrt.WriteHeader(c)
					wrt.Write([]byte(e.Error()))
					break
				}
			}
			return
		}
	}
	http.NotFound(wrt, req)
}
