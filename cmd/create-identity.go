package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/stellar/go-stellar-base/keypair"
)

// createIdentityCmd represents the add-identity command
var createIdentityCmd = &cobra.Command{
	Use:   "create-identity NAME",
	Short: "creates a new identity in the config at NAME",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			cmd.Usage()
			os.Exit(-1)
		}

		name := args[0]
		kp, err := keypair.Random()
		mustSucceed(err)
		mustSucceed(addIdentity(name, kp.Seed()))
		mustSucceed(saveConfig(cfgFile))

		fmt.Println("account id:", kp.Address())
		fmt.Println("secret key:", kp.Seed())
	},
}

func init() {
	RootCmd.AddCommand(createIdentityCmd)
}
