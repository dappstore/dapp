package dapp

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/jbenet/go-multihash"
	"github.com/pkg/errors"
)

// DefaultHorizonServer will be used to determine the manifest for a dapp when
// no custom servers are specified in a call to `Register`.
const DefaultHorizonServer = "https://horizon.stellar.org"

var dev = flag.Bool(
	"dapp.dev",
	false,
	"enables developer mode",
)

var printVersion = flag.Bool(
	"dapp.version",
	false,
	"print the current version and exit",
)

var printID = flag.Bool(
	"dapp.id",
	false,
	"print the app's id and exit",
)

var version = "devel"

type manifestDisagreementError struct{}

// ManifestHash resolves the multihash for the manifest of application `id`
// using `horizons`.
func ManifestHash(id string, horizons ...string) (multihash.Multihash, error) {
	var err error

	if len(horizons) == 0 {
		horizons = []string{DefaultHorizonServer}
	}
	manifestHashes := make([]multihash.Multihash, len(horizons))

	for i, server := range horizons {
		manifestHashes[i], err = loadSingleManifest(id, server)
		if err != nil {
			return nil, errors.Wrap(err, "load manifest hash failed")
		}
	}

	// TODO: ensure trust threshold is satisfied

	// return the first non-nil
	for _, h := range manifestHashes {
		if h != nil {
			return h, nil
		}
	}

	// TODO: use an error struct
	return nil, errors.New("could not load any manifest hashes")
}

// Register initializes the dapp system.
func Register(id string, horizons ...string) {
	if !flag.Parsed() {
		flag.Parse()
	}

	if *dev {
		return
	}

	if *printID {
		fmt.Println(id)
		os.Exit(0)
	}

	if *printVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	// load the "dapp-manifest" data field for the account at address from all the
	// horizon servers.  If they disagree, fatally error out (in the future,
	// perhaps retry).
	hash, err := ManifestHash(id, horizons...)
	if err != nil {
		errors.Print(err)
		os.Exit(1)
	}

	_ = hash
	log.Println(hash.B58String())

	// load the manifest from ipfs, parse it as json.
	// use semver to see if there is a
}

func loadSingleManifest(id string, server string) (multihash.Multihash, error) {
	url := fmt.Sprintf("%s/accounts/%s/data/dapp-manifest", server, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not create request")
	}

	req.Header.Add("Accept", "application/octet-stream")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "request errored")
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		// TODO: use a better error by interpetting the horizon problem response
		return nil, errors.New("request failed")
	}

	hash, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read failed")
	}

	return hash, nil
}
