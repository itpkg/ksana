package mux

import (
	"net/http"
	"regexp"
)

type Router struct {
	routes      []*Route
	middlewares []Middleware
}

func (p *Router) Get(pattern string, handlers ...Handler) {
	p.Add("GET", pattern, handlers...)
}

func (p *Router) Post(pattern string, handlers ...Handler) {
	p.Add("POST", pattern, handlers...)
}

func (p *Router) Put(pattern string, handlers ...Handler) {
	p.Add("PUT", pattern, handlers...)
}

func (p *Router) Patch(pattern string, handlers ...Handler) {
	p.Add("PATCH", pattern, handlers...)
}

func (p *Router) Delete(pattern string, handlers ...Handler) {
	p.Add("DELETE", pattern, handlers...)
}

func (p *Router) Add(method, pattern string, handlers ...Handler) {
	p.routes = append(
		p.routes,
		&Route{
			method:   method,
			pattern:  regexp.MustCompile(pattern),
			handlers: handlers,
		},
	)
}

func (p *Router) ServeHTTP(wrt http.ResponseWriter, req *http.Request) {
	for _, rt := range p.routes {
		if rt.Match(req) {
			hds := func(c *Context) (int, error) {
				for _, h := range rt.handlers {
					if s, e := h(c); e != nil {
						return s, e
					}
				}
				return 0, nil
			}

			for _, m := range rt.middlewares {
				hds = m(hds)
			}

			for _, m := range p.middlewares {
				hds = m(hds)
			}

			if c, e := hds(&Context{
				Writer:  wrt,
				Request: req,
				Params:  rt.Parse(req),
			}); e != nil {
				wrt.WriteHeader(c)
				wrt.Write([]byte(e.Error()))
			}
			return
		}
	}
	http.NotFound(wrt, req)
}
