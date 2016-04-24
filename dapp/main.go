package dapp

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/jbenet/go-multihash"
	"github.com/pkg/errors"
	"github.com/stellar/go-stellar-base/horizon"
	"github.com/stellar/go-stellar-base/keypair"
)

// App represents the identity for an application that is deployed using dapp
type App struct {
	ID  string
	Err error

	policies            []Policy
	verificationServers []string
	horizons            map[string]struct{}
}

// DirModifier can change a directory
type DirModifier func(path string) error

// Identity represents a single identity in the dapp system
type Identity struct {
	keypair.KP
}

// Policy values represent a policy that can change state on the app
type Policy interface {
	ApplyDappPolicy(*App) error
}

type VerifyAgainstAny struct{}
type VerifyAgainstAll struct{}

// Add ensures `path` is in ipfs
func Add(path string) (multihash.Multihash, error) {
	return ContentHash(path)
}

// ContentHash returns the content hash, according to ipfs, for `path`
func ContentHash(path string) (ret multihash.Multihash, err error) {
	stdout, err := exec.Command("ipfs", "add", "-q", path).Output()
	if err != nil {
		err = errors.Wrap(err, "ipfs add failed")
		return
	}

	hashes := strings.Split(strings.TrimSpace(string(stdout)), "\n")
	lastHash := hashes[len(hashes)-1]

	ret, err = multihash.FromB58String(lastHash)
	if err != nil {
		err = errors.Wrap(err, "failed decoding ipfs output")
		return
	}

	return
}

// CurrentUser returns the current process' identity within `app`
func CurrentUser(app string) *Identity {
	return loginSessions[app]
}

// Fund funds id on the stellar network using the configured friendbot.
func Fund(id *Identity) error {
	exists, err := identityExists(id)
	if err != nil {
		return errors.Wrap(err, "identity existence check errored")
	}

	if exists {
		// TODO: use an actual error struct, embed the network passphrase of the
		// horizon server consulted.
		return errors.New("identity already funded")
	}

	url := fmt.Sprintf("%s/friendbot?addr=%s", defaultHorizon, id.Address())

	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "friendbot error")
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		// TODO: use a better error by interpetting the horizon problem response
		return errors.New("friendbot failed")
	}

	return nil
}

// GetApplication initializes the dapp system for integrating applications.  It
// enables self-updateing and self-verifying features.
func GetApplication(id string, policies ...Policy) *App {
	app := &App{ID: id}
	err := app.InitializePolicies(policies)
	if err != nil {
		errors.Print(err)
		os.Exit(1)
	}

	if *printID {
		fmt.Println(id)
		os.Exit(0)
	}

	if *printVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	return app
}

// Login logs the current process into `app` as `user`, replacing any current
// session.
func Login(app string, user *Identity) {

	loginSessions[app] = user
}

// Logout logs the current process out of `app`
func Logout(app string) {
	delete(loginSessions, app)
}

// ManifestHash resolves the multihash for the manifest of application `id`
// using `horizons`.
func ManifestHash(id string, horizons ...string) (multihash.Multihash, error) {
	var err error

	if len(horizons) == 0 {
		return multihash.Multihash(""),
			errors.New("no verification servers specified")
	}
	manifestHashes := make([]multihash.Multihash, len(horizons))

	for i, server := range horizons {
		manifestHashes[i], err = loadRelease(server, id)
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

// NewIdentity creates and validates a new identity
func NewIdentity(seedOrAddress string) (*Identity, error) {
	kp, err := keypair.Parse(seedOrAddress)
	if err != nil {
		return nil, errors.Wrap(err, "parse identity")
	}

	return &Identity{kp}, nil
}

// SetDefaultHorizon sets the default horizon server
func SetDefaultHorizon(addy string) {
	defaultHorizon = addy
}

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
var defaultHorizon = horizon.DefaultTestNetClient.URL
var loginSessions map[string]*Identity

type manifestDisagreementError struct{}

func childDir(base multihash.Multihash, dir string) string {
	return fmt.Sprintf("ipfs/%s/%s/", base.B58String(), dir)
}

func init() {
	loginSessions = map[string]*Identity{}
}

func identityExists(id *Identity) (bool, error) {
	url := fmt.Sprintf("%s/accounts/%s", defaultHorizon, id.Address())

	resp, err := http.Get(url)
	if err != nil {
		return false, errors.Wrap(err, "load account data failed")
	}

	return (resp.StatusCode >= 200 && resp.StatusCode < 300), nil
}

func loadDirInto(base multihash.Multihash, dir string, local string) error {
	ipfsPath := childDir(base, dir)
	err := exec.Command("ipfs", "get", "-o", local, ipfsPath).Run()
	if err != nil {
		return errors.Wrap(err, "ipfs get failed")

	}

	return nil
}

func loadRelease(id string, server string) (multihash.Multihash, error) {
	hash, err := loadIdentityData(server, id, "dapp-release")
	if err != nil {
		return nil, errors.Wrap(err, "read identity data failed")
	}

	return hash, nil
}

func loadIdentityData(server, id, key string) ([]byte, error) {
	url := fmt.Sprintf("%s/accounts/%s/data/%s", server, id, key)
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

func loadIdentityMultihash(
	server,
	id,
	key string,
) (ret multihash.Multihash, err error) {

	ret, err = loadIdentityData(server, id, key)
	if err != nil {
		err = errors.Wrap(err, "read identity data failed")
		return
	}

	return
}

func verifyDir(base multihash.Multihash, dir string) error {
	ipfsPath := childDir(base, dir)
	_, err := exec.Command("ipfs", "ls", ipfsPath).Output()
	if err != nil {
		return errors.Wrap(err, "ipfs ls failed")
	}

	return nil
}

func verifyPublication(server, id, path string) (bool, error) {
	publishedHash, err := loadIdentityMultihash(server, id, "dapp-publications")
	if err != nil {
		return false, errors.Wrap(err, "get publication hash failed")
	}

	err = verifyDir(publishedHash, id)
	if err != nil {
		return false, errors.Wrap(err, "directory verification failed")
	}

	//TODO: load the manifest, verify signatures against binaries

	return true, nil
}

// ModifyDir loads an ipfs dir, modifies it according to `fn` and
// commits it back to ipfs, returning the hash
func ModifyDir(
	base multihash.Multihash,
	dir string,
	fn DirModifier,
) (ret multihash.Multihash, err error) {
	err = verifyDir(base, dir)
	if err != nil {
		err = errors.Wrap(err, "initial dir verification failed")
		return
	}

	next, err := ioutil.TempDir("", "dapp-modify-dir")
	if err != nil {
		return
	}
	defer os.RemoveAll(dir)

	err = loadDirInto(base, dir, next)
	if err != nil {
		err = errors.Wrap(err, "failed to populate temp dir")
		return
	}

	err = fn(next)
	if err != nil {
		err = errors.Wrap(err, "modify callback errored")
		return
	}

	ret, err = Add(next)
	if err != nil {
		err = errors.Wrap(err, "ipfs add dailed")
		return
	}

	return
}
