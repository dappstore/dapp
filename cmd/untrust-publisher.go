package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// untrustPublisherCmd represents the add-identity command
var untrustPublisherCmd = &cobra.Command{
	Use:   "untrust-publisher ID_OR_ALIAS",
	Short: "removes trust in publisher",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			cmd.Usage()
			os.Exit(-1)
		}

		idOrAlias := args[0]
		id, err := resolveAlias(idOrAlias)
		mustSucceed(err)

		var next []string
		for i, addy := range config.TrustedPublishers {
			if addy == id.Address() {
				next = append(
					config.TrustedPublishers[:i],
					config.TrustedPublishers[i+1:]...,
				)
			}
		}

		if next == nil {
			fail("publisher not found", -1)
		}

		config.TrustedPublishers = next
		mustSucceed(saveConfig(cfgFile))
	},
}

func init() {
	RootCmd.AddCommand(untrustPublisherCmd)
	untrustPublisherCmd.Flags().StringVarP(&newPublisherAlias, "alias", "", "", "Also add alias for publisher")
}
