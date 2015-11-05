package mux

import (
	"fmt"
	"net/http"
)

type Middleware func(Handler) Handler

type Handler func(*Context) (int, error)

func Start(router *Router, port int) error {
	for _, en := range engines {
		en.Mount(router)
	}
	http.Handle("/", router)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
