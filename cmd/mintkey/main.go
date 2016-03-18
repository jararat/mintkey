package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {

	privKeyFlag := cli.StringFlag{
		Name:  "priv-key",
		Value: defaultPath("priv_key"),
		Usage: "Wire expression of message to sign",
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
