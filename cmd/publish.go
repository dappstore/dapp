package cmd

import (
	"fmt"

	"github.com/dappstore/go-dapp/protocols/dfs"
	"github.com/dappstore/go-dapp/protocols/publish"
	"github.com/spf13/cobra"
)

// var pathToPublish string

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish [FILES ...]",
	Short: "Publishes data after signing with the current dapp identity",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		login()
		pfs := dfs.New(app.Providers)
		ppublish := publish.New(app.Providers)

		hash, err := pfs.StoreLocalPaths(args)
		mustSucceed(err)

		// TODO: add publication file to `contents`

		tx, hash, err := ppublish.SetPublications(app.CurrentUser(), hash)
		mustSucceed(err)

		fmt.Println("publisher:", app.CurrentUser())
		fmt.Println("new publications hash:", hash)
		fmt.Println("transaction hash:", tx)
	},
}

func init() {
	RootCmd.AddCommand(publishCmd)
	// publishCmd.Flags().StringVarP(
	// 	&pathToPublish, "path", "",
	// 	"",
	// 	"path to publish FILES at",
	// )
}
