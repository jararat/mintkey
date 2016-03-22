package main

import (
	"fmt"

	"github.com/codegangsta/cli"
	. "github.com/tendermint/go-common"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-wire"
)

// cli entrypoint to show key information
func cmdShow(c *cli.Context) {

	// Load the privKey to sign with.
	privKeyPath := c.String("priv-key")
	privKey, err := loadPrivKey(privKeyPath)
	if err != nil {
		Exit("Error loading private key: " + err.Error())
	}

	// Show pubkey
	pubKey := privKey.PubKey()
	fmt.Printf("PubKey Address: %X\n", pubKey.Address())
	fmt.Printf("PubKey Bytes:   %X\n", pubKey.Bytes())
	fmt.Printf("PubKey JSON:    %v\n", string(wire.JSONBytesPretty(struct{ crypto.PubKey }{pubKey})))
}
