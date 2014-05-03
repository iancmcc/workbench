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
	VERSION = "0.1.0"
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
		fmt.Printf(banner, VERSION)
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
	jig.Build(jf, ctx.Args())
}

func Initialize(ctx *cli.Context) {
	var (
		pwd, name string
		err       error
	)
	args := ctx.Args()
	if args.Present() {
		name = args.First()
	} else {
		name = "build"
	}
	if pwd, err = os.Getwd(); err != nil {
		log.Critical("Unable to get current directory: %+v", err)
		os.Exit(1)
	}
	if err := jig.Initialize(pwd, name, ctx.String("image"),
		ctx.StringSlice("pre"),
		ctx.StringSlice("build"),
		ctx.StringSlice("post")); err != nil {
		log.Error("Initialization of Jigfile failed: %+v", err)
		os.Exit(1)
	}
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
	initialize := cli.Command{
		Name:   "init",
		Usage:  "Initialize a Jigfile in the current directory",
		Action: Initialize,
	}
	initialize.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "image, i",
			Value: "jigs/base",
			Usage: "Base image to use for the build",
		},
		cli.StringSliceFlag{
			Name:  "pre",
			Value: &cli.StringSlice{},
			Usage: "Commands to be executed as root before the build",
		},
		cli.StringSliceFlag{
			Name:  "build, b",
			Value: &cli.StringSlice{},
			Usage: "Build commands to be executed",
		},
		cli.StringSliceFlag{
			Name:  "post",
			Value: &cli.StringSlice{},
			Usage: "Commands to be executed as root after the build",
		},
	}
	build.Flags = []cli.Flag{}
	app.Commands = []cli.Command{build, initialize}
	app.Version = VERSION
	app.Run(os.Args)
}
