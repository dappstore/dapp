package cmd

import (
	"fmt"
	"os"

	"github.com/dappstore/go-dapp/stellar"
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

		fmt.Println("creating identity...")
		id, err := app.Providers.RandomIdentity()
		mustSucceed(err)

		kp := id.(*stellar.Identity).KP.(*keypair.Full)
		fmt.Println("saving local config...")
		mustSucceed(addIdentity(name, kp.Seed()))
		mustSucceed(saveConfig(cfgFile))

		// Add "Identity"
		fmt.Println("publicly accouncing identity...")
		hash, err := app.Providers.AnnounceIdentity(id)
		mustSucceed(err)

		fmt.Println("new identity successfully announced in tx:", hash)
		fmt.Println("public key:", kp.Address())
		fmt.Println("secret key:", kp.Seed())
	},
}

func init() {
	RootCmd.AddCommand(createIdentityCmd)
}
