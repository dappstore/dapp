package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var newPublisherAlias string

// trustPublisherCmd represents the add-identity command
var trustPublisherCmd = &cobra.Command{
	Use:   "trust-publisher ID_OR_ALIAS",
	Short: "records trust in publisher",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			cmd.Usage()
			os.Exit(-1)
		}

		idOrAlias := args[0]
		id, err := resolveAlias(idOrAlias)
		mustSucceed(err)

		newAddy := id.PublicKey()

		for _, existing := range config.TrustedPublishers {
			if newAddy == existing {
				fail("publisher is already trusted", -1)
			}
		}

		config.TrustedPublishers = append(config.TrustedPublishers, newAddy)

		if newPublisherAlias != "" {
			mustSucceed(addIdentity(newPublisherAlias, newAddy))
		}

		mustSucceed(saveConfig(cfgFile))
	},
}

func init() {
	RootCmd.AddCommand(trustPublisherCmd)
	trustPublisherCmd.Flags().StringVarP(&newPublisherAlias, "alias", "", "", "Also add alias for publisher")
}
