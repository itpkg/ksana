package main

import (
	_ "github.com/lib/pq"

	"github.com/itpkg/ksana"
	_ "github.com/itpkg/ksana/orm"
)

func main() {
	if err := ksana.Run(); err != nil {
		panic(err)
	}
}
