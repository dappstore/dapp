package cmd

import (
	"fmt"
	"os"

	"github.com/nullstyle/dapp"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var cfgFile string
var identity string
var app *dapp.App

func getIdentity(alias string) *dapp.Identity {
	seedOrAddress, ok := viper.GetStringMapString("Identities")[alias]
	if !ok {
		fmt.Printf("no identity: %s", identity)
		os.Exit(-1)
	}

	id, err := dapp.NewIdentity(seedOrAddress)
	if err != nil {
		errors.Print(err)
		os.Exit(-1)
	}

	return id
}
