package mux

import (
	"fmt"
	"net/http"
)

type Middleware func(Handler) Handler

type Handler func(*Context) (int, error)

func Start(router *Router, port int) error {
	http.Handle("/", router)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
