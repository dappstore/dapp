package dapp

import (
	"log"
	"os"

	"github.com/pkg/errors"
)

// ApplyPolicy applies `p` t `a`
func (a *App) ApplyPolicy(p Policy) error {
	err := p.ApplyDappPolicy(a)
	if err != nil {
		return errors.Wrap(err, "failed applying policy")
	}

	return nil
}

// InitializePolicies applies `policies` t `a`
func (a *App) InitializePolicies(policies []Policy) error {
	a.policies = append(a.policies, policies...)
	a.do(
		func() error { return a.ApplyPolicy(&VerifySelf{Publisher: a.ID}) },
	)

	for _, p := range policies {
		a.do(
			func() error { return a.ApplyPolicy(p) },
		)
	}

	if len(a.verificationServers) == 0 {
		a.do(
			func() error { return a.ApplyPolicy(&VerifyUsing{defaultHorizon}) },
		)
	}

	// Post policies
	// a.do(
	// 	func() error { return a.ApplyPolicy(&Verify{}) },
	// )

	return nil
}

// Fund funds `user`
func (a *App) Fund(id *Identity) error {
	return Fund(id)
}

// Login logs `user` into `a`
func (a *App) Login(user *Identity) {
	Login(a.ID, user)
}

// do is a helper to only perform actions while an app remains "un-errored".
func (a *App) do(fns ...func() error) error {
	if a.Err != nil {
		return a.Err
	}

	for _, fn := range fns {
		a.Err = fn()
		if a.Err != nil {
			return a.Err
		}
	}
	return nil
}

func (a *App) verifyPublished() error {
	// load the "dapp-manifest" data field for the account at address from all the
	// horizon servers.  If they disagree, fatally error out (in the future,
	// perhaps retry).
	hash, err := ManifestHash(a.ID, a.verificationServers...)
	if err != nil {
		errors.Print(err)
		os.Exit(1)
	}

	_ = hash
	log.Println(hash.B58String())
	return nil
}
