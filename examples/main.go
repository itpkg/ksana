package main

import (
	"github.com/itpkg/ksana"
	_ "github.com/itpkg/ksana/base"
	_ "github.com/lib/pq"
)

func main() {
	if err := ksana.Run(); err != nil {
		panic(err)
	}
}
