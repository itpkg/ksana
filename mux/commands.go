package mux

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/itpkg/ksana/cmd"
)

func init() {
	cmd.Register(cli.Command{
		Name:    "routes",
		Aliases: []string{"ro"},
		Usage:   "print out all defined routes in match order.",
		Action: func(c *cli.Context) {
			rt := New()
			for _, en := range engines {
				en.Mount(rt)
			}
			for _, r := range rt.routes {
				fmt.Fprintf(os.Stdout, "%s %v \n", r.method, r.pattern)
			}
		},
	})
}
