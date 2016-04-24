package cmd

import (
	"github.com/spf13/cobra"
)

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish [file ...]",
	Short: "Publishes data after signing with the current dapp identity",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO:
		// - add files to ipfs
		// - Create directory with publication file, all published files at current  identity
		// - Add publication to ipfs
	},
}

func init() {
	RootCmd.AddCommand(publishCmd)
}
