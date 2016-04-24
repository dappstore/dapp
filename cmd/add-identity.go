package cmd

import (
	"os"

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

		mustSucceed(addIdentity(name, seedOrAddress))
		mustSucceed(saveConfig(cfgFile))
	},
}

func init() {
	RootCmd.AddCommand(addIdentityCmd)
}
