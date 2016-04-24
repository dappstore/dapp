package dapp

// VerifyUsing represents the dapp policy that when applied to the current
// process causes dapp consider using `Horizon` when performing the verification
// protocol.  Other policies may augment how the server is concretely used.  See
// `VerifyAgainstAll` and `VerifyAgainstAny`.
type VerifyUsing struct {
	Horizon string
}

// ApplyDappPolicy applies `p` to `app`
func (p *VerifyUsing) ApplyDappPolicy(app *App) error {
	return nil
}

var _ Policy = &VerifyUsing{}
