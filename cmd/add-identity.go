package cmd

import (
	"fmt"
	"os"

	"github.com/dappstore/dapp/dapp"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// addIdentityCmd represents the add-identity command
var addIdentityCmd = &cobra.Command{
	Use:   "add-identity NAME ID",
	Short: "records an identity in the config",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 2 {
			cmd.Usage()
			os.Exit(-1)
		}

		name := args[0]
		seedOrAddress := args[1]

		if config.Identities[name] != "" {
			message := fmt.Sprintf("identity already exists")
			fail(errors.New(message), -1)
		}

		_, err := dapp.NewIdentity(seedOrAddress)
		if err != nil {
			message := fmt.Sprintf("ID is invalid")
			fail(errors.New(message), -1)
		}

		config.Identities[name] = seedOrAddress

		mustSucceed(saveConfig(cfgFile))
	},
}

func init() {
	RootCmd.AddCommand(addIdentityCmd)
}
