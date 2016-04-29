package cmd

import (
	"fmt"
	// "log"
	"path/filepath"

	"github.com/dappstore/go-dapp"
	"github.com/dappstore/go-dapp/protocols/dfs"
	"github.com/dappstore/go-dapp/protocols/publish"
	"github.com/gosuri/uilive"
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
		writer := uilive.New()
		writer.Start()

		pdfs := dfs.New(app.Providers)
		ppublish := publish.New(app.Providers, app.Providers)

		fmt.Fprintln(writer, "loading current publications hash...")
		source, err := ppublish.GetPublications(app.CurrentUser())
		mustSucceed(err)

		current := source

		for _, path := range args {
			var toAdd dapp.Hash
			name := filepath.Base(path)
			toAdd, err = pdfs.StorePath(path)
			mustSucceed(err)
			fmt.Fprintf(writer, "merging %s into %s\n", path, current)
			current, err = pdfs.MergeAtPath(current, name, toAdd)
			mustSucceed(err)
		}

		fmt.Fprintf(writer, "setting %s as publications...\n", current)
		tx, pubs, err := ppublish.SetPublications(app.CurrentUser(), current)
		mustSucceed(err)

		fmt.Fprint(writer, "Success!\n\n")
		writer.Stop()

		fmt.Println("publisher:", app.CurrentUser())
		fmt.Println("content hash:", current)
		fmt.Println("publications hash:", pubs)
		fmt.Println("transaction hash:", tx)
	},
}

func init() {
	RootCmd.AddCommand(publishCmd)
}
