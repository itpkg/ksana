package ksana

import (
	"os"

	"github.com/itpkg/ksana/cmd"
)

func Run() error {
	app := cmd.New(
		"ksana",
		"A web framwork for go like rails",
		"v20151104")
	return app.Run(os.Args)
}
