package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dappstore/go-dapp/protocols/dfs"
	"github.com/dappstore/go-dapp/protocols/publish"
	"github.com/spf13/cobra"
)

// getPublicationCmd represents the publish command
var getPublicationCmd = &cobra.Command{
	Use:   "get-publication ID_OR_ALIAS PATH",
	Short: "Loads a publication into a local sub directory",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 2 {
			cmd.Usage()
			os.Exit(-1)
		}

		idOrAlias := args[0]
		path := args[1]

		id, err := resolveAlias(idOrAlias)
		mustSucceed(err)
		fmt.Println("resolved publisher:", id)

		// TODO: check publisher trust

		pfs := dfs.New(app.Providers)
		ppublish := publish.New(app.Providers, app.Providers)

		contents, err := ppublish.GetPublications(id)
		mustSucceed(err)
		fmt.Println("current publications hash:", contents)

		// TODO: verify claims before loading data

		name := filepath.Base(path)
		dir, err := pfs.LoadTemp(contents)
		localPath := filepath.Join(dir, path)
		err = os.Rename(localPath, name)

		if err != nil {
			fail("failed to move local copy of publication to current dir", -1)
		}
	},
}

func init() {
	RootCmd.AddCommand(getPublicationCmd)
}
