package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/dappstore/go-dapp"
	apps "github.com/dappstore/go-dapp/app"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var cfgFile string
var identity string
var app *apps.App
var config struct {
	CacheDir          string
	Identities        map[string]string
	TrustedPublishers []string
}

func addIdentity(name, seedOrAddress string) error {
	if config.Identities[name] != "" {
		return errors.New("identity already exists")
	}

	_, err := app.Providers.ParseIdentity(seedOrAddress)
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
		fail(fmt.Sprintf("no identity: %s", identity), -1)
	}

	id, err := app.Providers.ParseIdentity(seedOrAddress)
	mustSucceed(err)

	return id
}

func login() {
	id := getIdentity(identity)
	app.Login(id)
}

func resolveAlias(idOrAlias string) (ret dapp.Identity, err error) {
	ret, err = app.Providers.ParseIdentity(idOrAlias)
	if err == nil {
		return
	}

	aliasID := config.Identities[idOrAlias]
	if aliasID == "" {
		err = errors.New("resolveAlias: no alias found")
		return
	}

	ret, err = app.Providers.ParseIdentity(aliasID)
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
