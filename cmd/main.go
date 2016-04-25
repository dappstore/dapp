package cmd

import (
	"bytes"
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
	CacheDir          string
	Identities        map[string]string
	TrustedPublishers []string
}

func addIdentity(name, seedOrAddress string) error {
	if config.Identities[name] != "" {
		return errors.New("identity already exists")
	}

	_, err := dapp.NewIdentity(seedOrAddress)
	if err != nil {
		return errors.Wrap(err, "addIdentity: invalid identity")
	}

	config.Identities[name] = seedOrAddress
	return nil
}

func fail(err string, code int) {
	if err == "" {
		panic("cannot fail with empty message")
	}

	fmt.Fprintln(os.Stderr, err)
	os.Exit(code)
}

func mustSucceed(err error) {
	if err != nil {
		var buf bytes.Buffer
		errors.Fprint(&buf, err)
		fail(buf.String(), -1)
	}
}

func getIdentity(alias string) dapp.Identity {
	seedOrAddress, ok := viper.GetStringMapString("Identities")[alias]
	if !ok {
		fmt.Printf("no identity: %s", identity)
		os.Exit(-1)
	}

	id, err := dapp.NewIdentity(seedOrAddress)
	mustSucceed(err)

	return id
}

func resolveAlias(idOrAlias string) (ret dapp.Identity, err error) {
	ret, err = dapp.NewIdentity(idOrAlias)
	if err == nil {
		return
	}

	aliasID := config.Identities[idOrAlias]
	if aliasID == "" {
		err = errors.New("resolveAlias: no alias found")
		return
	}

	ret, err = dapp.NewIdentity(aliasID)
	if err != nil {
		err = errors.Wrap(err, "resolveAlias: corrupt alias")
		return
	}

	return
}

func saveConfig(path string) error {

	toSave := map[string]interface{}{
		"CacheDir":          config.CacheDir,
		"Identities":        config.Identities,
		"TrustedPublishers": config.TrustedPublishers,
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
