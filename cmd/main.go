package cmd

import (
	"fmt"
	"os"

	"github.com/dappstore/dapp/dapp"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var cfgFile string
var identity string
var app *dapp.App
var config struct {
	CacheDir   string
	Identities map[string]string
}

func fail(err error, code int) {
	if err == nil {
		panic("cannot fail with nil error")
	}

	errors.Print(err)
	os.Exit(code)
}

func mustSucceed(err error) {
	if err != nil {
		fail(err, -1)
	}
}

func getIdentity(alias string) *dapp.Identity {
	seedOrAddress, ok := viper.GetStringMapString("Identities")[alias]
	if !ok {
		fmt.Printf("no identity: %s", identity)
		os.Exit(-1)
	}

	id, err := dapp.NewIdentity(seedOrAddress)
	mustSucceed(err)

	return id
}

func saveConfig(path string) error {

	toSave := map[string]interface{}{
		"CacheDir":   config.CacheDir,
		"Identities": config.Identities,
	}

	b, err := yaml.Marshal(toSave)
	if err != nil {
		return errors.Wrap(err, "save-config: marshal to yaml failed")
	}

	f, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "save-config: failed to create file")
	}

	defer f.Close()

	f.WriteString(string(b))

	return nil
}
