package cmd

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// fundCmd represents the fund command
var fundCmd = &cobra.Command{
	Use:   "fund [identity]",
	Short: "funds a identity",
	Long:  `The fund command registers an identity with the stellar network, which can then be used to publish data as part of the dapp protocols.`,
	Run: func(cmd *cobra.Command, args []string) {

		id := getIdentity(args[0])

		err := app.Fund(id)
		if err != nil {
			errors.Print(err)
			os.Exit(-1)
		}
	},
}

func init() {
	RootCmd.AddCommand(fundCmd)
}
