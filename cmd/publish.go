package cmd

import (
	"path/filepath"

	"github.com/dappstore/dapp/dapp"
	"github.com/spf13/cobra"
)

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish [file ...]",
	Short: "Publishes data after signing with the current dapp identity",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {

		contents := map[string]dapp.Hash{}

		// Add all paths to ipfs
		for _, path := range args {
			var err error
			name := filepath.Base(path)
			contents[name], err = app.StorePath(path)
			mustSucceed(err)
		}

		// TODO: add publication file to `contents`

		// pubHash, err := app.StoreMap(contents)
		// mustSucceed(err)

		// TODO:
		// - Create directory with publication file, all published files at current  identity
		// - Add publication to ipfs
	},
}

func init() {
	RootCmd.AddCommand(publishCmd)
}
