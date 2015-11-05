package mux

import (
	"net/http"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Params  map[string]interface{}
}

func (p *Context) Json(val interface{}) {
	//todo
}

func (p *Context) Xml(val interface{}) {
	//todo
}
