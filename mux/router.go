package mux

import (
	"net/http"
	"regexp"
)

type Router struct {
	routes      []*Route
	middlewares []Middleware
}

func (p *Router) Use(ms ...Middleware) {
	p.middlewares = append(p.middlewares, ms...)
}

func (p *Router) Get(pattern string, middlewares []Middleware, handlers ...Handler) {
	p.Add("GET", pattern, middlewares, handlers...)
}

func (p *Router) Post(pattern string, middlewares []Middleware, handlers ...Handler) {
	p.Add("POST", pattern, middlewares, handlers...)
}

func (p *Router) Put(pattern string, middlewares []Middleware, handlers ...Handler) {
	p.Add("PUT", pattern, middlewares, handlers...)
}

func (p *Router) Patch(pattern string, middlewares []Middleware, handlers ...Handler) {
	p.Add("PATCH", pattern, middlewares, handlers...)
}

func (p *Router) Delete(pattern string, middlewares []Middleware, handlers ...Handler) {
	p.Add("DELETE", pattern, middlewares, handlers...)
}

func (p *Router) Add(method, pattern string, middlewares []Middleware, handlers ...Handler) {
	p.routes = append(
		p.routes,
		&Route{
			method:      method,
			pattern:     regexp.MustCompile(pattern),
			handlers:    handlers,
			middlewares: middlewares,
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

			for i := len(rt.middlewares) - 1; i >= 0; i-- {
				m := rt.middlewares[i]
				hds = m(hds)
			}

			for i := len(p.middlewares) - 1; i >= 0; i-- {
				m := p.middlewares[i]
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

//==============================================================================

func New() *Router {
	return &Router{
		middlewares: make([]Middleware, 0),
		routes:      make([]*Route, 0),
	}
}
