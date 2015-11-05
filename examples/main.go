package main

import (
	_ "github.com/lib/pq"

	kc "github.com/itpkg/ksana/cmd"
	_ "github.com/itpkg/ksana/orm"
)

func main() {
	kc.Run()
}
