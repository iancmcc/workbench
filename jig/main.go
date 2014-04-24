package main

import (
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/iancmcc/jig"
)

func Build(c *cli.Context) {
	jig.ParseJigfile(strings.NewReader(""))
}

func main() {
	app := cli.NewApp()
	app.Name = "jig"
	app.Usage = "Portable build system"
	app.Commands = []cli.Command{
		{
			Name:   "build",
			Usage:  "Build artifacts from a Jigfile",
			Action: Build,
		},
	}
	app.Run(os.Args)
}
