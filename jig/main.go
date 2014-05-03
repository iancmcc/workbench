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

func Build(ctx *cli.Context) {
	var (
		jf  *jig.Jigfile
		err error
		pwd string
	)
	if !ctx.GlobalBool("quiet") {
		fmt.Printf(banner, version)
		fmt.Println("Running builds...")
	}
	if pwd, err = os.Getwd(); err != nil {
		log.Critical("Unable to get current directory")
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

func jigInit(ctx *cli.Context) error {
	format := logging.MustStringFormatter("[%{level}] %{message}")
	logBackend := logging.NewLogBackend(os.Stdout, "", stdlog.LstdFlags)
	if isTTY(os.Stdout) {
		logBackend.Color = true
	}
	logging.SetFormatter(format)
	logging.SetBackend(logBackend)
	logging.SetLevel(logging.Level(ctx.GlobalInt("level")), "jig")
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "jig"
	app.Usage = "Portable build system"
	app.Before = jigInit
	app.Flags = []cli.Flag{
		cli.IntFlag{"level, l", 4, "log level (5 for debug, 0 for silent)"},
		cli.BoolFlag{"quiet, q", "don't display banner or messages"},
	}
	build := cli.Command{
		Name:   "build",
		Usage:  "Build artifacts from a Jigfile",
		Action: Build,
	}
	build.Flags = []cli.Flag{
		cli.StringFlag{"output, o", "./output", "output directory for artifacts"},
	}
	app.Commands = []cli.Command{build}
	app.Run(os.Args)
}
