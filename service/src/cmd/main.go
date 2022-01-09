package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/judwhite/go-svc"
	"github.com/urfave/cli/v2"
)

var app = &cli.App{
	Version: fmt.Sprintf("%s(%s)", version, tag),
	Name:    "sms_blog",
	Usage:   "sms blog service",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "debug",
			Aliases: []string{"d"},
			Usage:   "开启debug",
		},
		&cli.StringFlag{
			Name:    "address",
			Aliases: []string{"a"},
			Value:   ":8088",
			Usage:   "http listen address",
		},
		&cli.StringFlag{
			Name:    "conf",
			Aliases: []string{"c"},
			Value:   "./config.json",
			Usage:   "config file path",
		},
		&cli.StringFlag{
			Name:        "log-level",
			DefaultText: "info",
			Usage:       "debug|info|warn|error",
		},
	},
	Action: func(c *cli.Context) error {
		prg := &program{ctx: c}
		return svc.Run(prg, syscall.SIGINT, syscall.SIGTERM)
	},
}

func main() {
	app.Run(os.Args)
}
