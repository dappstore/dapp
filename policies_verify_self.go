package dapp

// VerifySelf is a policy that causes the binary to verify itself as an
// installation of the application published by `Publisher`, according to the
// dapp publisher protocol.
type VerifySelf struct {
	Publisher string
}

// ApplyDappPolicy applies `p` to `app`
func (p *VerifySelf) ApplyDappPolicy(app *App) error {
	return nil
}

var _ Policy = &VerifySelf{}
