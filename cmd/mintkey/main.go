package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {

	privKeyFlag := cli.StringFlag{
		Name:  "priv-key",
		Value: defaultPath("priv_key"),
		Usage: "path to priv-key",
	}
	showPrivFlag := cli.BoolFlag{
		Name:  "show-priv",
		Usage: "whether to show private information",
	}
	/*
		quietFlag := cli.BoolFlag{
			Name:  "quiet",
			Value: false,
			Usage: "Without debug output",
		}
	*/

	app := cli.NewApp()
	app.Name = "mintkey"
	app.Usage = "Sign stuff"
	app.Commands = []cli.Command{
		{
			Name:  "gen",
			Usage: "generate a new private key",
			Action: func(c *cli.Context) {
				cmdGen(c)
			},
		},
		{
			Name:  "show",
			Usage: "show key information",
			Flags: []cli.Flag{privKeyFlag, showPrivFlag},
			Action: func(c *cli.Context) {
				cmdShow(c)
			},
		},
		{
			Name:  "sign",
			Usage: "sign some bytes",
			Flags: []cli.Flag{privKeyFlag},
			Action: func(c *cli.Context) {
				cmdSign(c)
			},
		},
	}

	app.Run(os.Args)
}
