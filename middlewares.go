package ksana

import (
	"github.com/itpkg/ksana/logging"
	"github.com/itpkg/ksana/mux"
)

func Logger(log logging.Logger) mux.Middleware {
	return func(next mux.Handler) mux.Handler {
		return func(ctx *mux.Context) (int, error) {
			log.Info("%s %s", ctx.Request.Method, ctx.Request.URL.Path)
			code, err := next(ctx)
			return code, err
		}
	}
}
