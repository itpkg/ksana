package main

import (
	"log"

	"github.com/itpkg/ksana/mux"
)

func M11() mux.Middleware {
	return func(next mux.Handler) mux.Handler {
		return func(ctx *mux.Context) (int, error) {
			log.Println("M11 BEGIN")
			code, err := next(ctx)
			log.Println("M11 END")
			return code, err
		}
	}
}

func M13() mux.Middleware {
	return func(next mux.Handler) mux.Handler {
		return func(ctx *mux.Context) (int, error) {
			log.Println("M13 BEGIN")
			code, err := next(ctx)
			log.Println("M13 END")
			return code, err
		}
	}
}

func M12() mux.Middleware {
	return func(next mux.Handler) mux.Handler {
		return func(ctx *mux.Context) (int, error) {
			log.Println("M12 BEGIN")
			code, err := next(ctx)
			log.Println("M12 END")
			return code, err
		}
	}
}

func M21() mux.Middleware {
	return func(next mux.Handler) mux.Handler {
		return func(ctx *mux.Context) (int, error) {
			log.Println("M21 BEGIN")
			code, err := next(ctx)
			log.Println("M21 END")
			return code, err
		}
	}
}

func M22() mux.Middleware {
	return func(next mux.Handler) mux.Handler {
		return func(ctx *mux.Context) (int, error) {
			log.Println("M22 BEGIN")
			code, err := next(ctx)
			log.Println("M22 END")
			return code, err
		}
	}
}

func M23() mux.Middleware {
	return func(next mux.Handler) mux.Handler {
		return func(ctx *mux.Context) (int, error) {
			log.Println("M23 BEGIN")
			code, err := next(ctx)
			log.Println("M23 END")
			return code, err
		}
	}
}

func H1(_ *mux.Context) (int, error) {
	log.Println("RUN H1")
	return 0, nil
}

func H2(_ *mux.Context) (int, error) {
	log.Println("RUN H2")
	return 0, nil
}

func H3(_ *mux.Context) (int, error) {
	log.Println("RUN H3")
	return 0, nil
}

func main() {
	rt := mux.New()
	rt.Use(M11(), M12(), M13())

	rt.Get("/test", []mux.Middleware{M21(), M22(), M23()}, H1, H2, H3)

	if err := mux.Start(rt, 8080); err != nil {
		panic(err)
	}
}
