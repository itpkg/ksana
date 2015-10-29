package mux

import (
	"fmt"
	"net/http"
)

type HttpHandler func(c *Context) (int, error)

func Start(router *Router, port int) error {
	http.Handle("/", router)
	return http.ListenAndServe(fmt.Sprintf(":%d", port, router), nil)
}
