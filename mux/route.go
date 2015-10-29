package mux

type Route struct {
	Method   string
	Pattern  string
	Handlers []HttpHandler
}
