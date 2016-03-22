package main

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
	. "github.com/tendermint/go-common"
	"github.com/tendermint/go-wire/expr"
)

// cli entrypoint to sign messages
func cmdSign(c *cli.Context) {

	// Load the privKey to sign with.
	privKeyPath := c.String("priv-key")
	privKey, err := loadPrivKey(privKeyPath)
	if err != nil {
		Exit("Error loading private key: " + err.Error())
	}

	// Get signbytes
	signExprStr := c.Args()[0]
	signExpr, err := expr.ParseReader("", strings.NewReader(signExprStr))
	if err != nil {
		Exit("Error parsing input: " + err.Error())
	}
	signBytes, err := signExpr.(expr.Byteful).Bytes()
	if err != nil {
		Exit("Error serializing parsed input: " + err.Error())
	}

	fmt.Printf("Wire expression: %v\n", signExpr)
	fmt.Printf("Wire bytes: %X\n", signBytes)
	signature := privKey.Sign(signBytes)
	fmt.Printf("Signature: %X\n", signature.Bytes())
}
