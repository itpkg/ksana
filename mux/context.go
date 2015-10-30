package mux

import (
	"net/http"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	params  map[string]interface{}
}

func (p *Context) Set(key string, val interface{}) {
	p.params[key] = val
}

func (p *Context) Get(key string) interface{} {
	return p.params[key]
}
