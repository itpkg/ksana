package mux

import (
	"net/http"
)

type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
	Env      map[string]interface{}
}
