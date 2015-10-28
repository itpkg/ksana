package main

import (
	_ "github.com/lib/pq"

	ks "github.com/itpkg/ksana"
	_ "github.com/itpkg/ksana/base"
)

func main() {
	ks.Run()
}
