package cmd

import (
	"fmt"
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

		tx, hash, err := app.PublishMap(contents)
		mustSucceed(err)

		fmt.Println("publisher:", app.CurrentUser())
		fmt.Println("new publications hash:", hash)
		fmt.Println("transaction hash:", tx)
	},
}

func init() {
	RootCmd.AddCommand(publishCmd)
}
