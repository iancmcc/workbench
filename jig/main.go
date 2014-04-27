package main

import (
	"fmt"
	stdlog "log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/iancmcc/jig"
	"github.com/op/go-logging"
)

var (
	log     = logging.MustGetLogger("jig")
	version = "0.0.1"
)

const banner string = `      _ _
     (_|_)__ _      v%s
     | | / _` + "`" + ` |     Portable builds
    _/ |_\__, |     for everyone!
   |__/  |___/

`

func Build(c *cli.Context) {
	var (
		jf  *jig.Jigfile
		err error
		pwd string
	)
	fmt.Printf(banner, version)
	if pwd, err = os.Getwd(); err != nil {
		log.Critical("Unable to do something")
		fmt.Println(err)
		os.Exit(1)
	}
	log.Debug(fmt.Sprintf("Looking for Jigfile in %s", pwd))
	if jf, err = jig.ParseJigfilePath(pwd); err != nil {
		log.Critical("No Jigfile found.")
		log.Debug(fmt.Sprintf("%v", err))
		os.Exit(1)
	}
	jig.Build(jf)
}

func main() {
	format := logging.MustStringFormatter("[%{level}] %{message}")
	logBackend := logging.NewLogBackend(os.Stdout, "", stdlog.LstdFlags)
	if jig.IsTTY(os.Stdout) {
		logBackend.Color = true
	}
	logging.SetFormatter(format)
	logging.SetBackend(logBackend)
	logging.SetLevel(logging.DEBUG, "jig")

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
