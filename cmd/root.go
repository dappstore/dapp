package cmd

import (
	"fmt"
	"os"

	"github.com/dappstore/dapp/dapp"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "dapp",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dapp.yaml)")
	RootCmd.PersistentFlags().StringVar(&identity, "identity", "default", "the identity used to authorize the command")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/dapp")
	viper.AddConfigPath("$HOME/.dapp")
	viper.SetDefault("CacheDir", os.ExpandEnv("$HOME/.dapp/cache"))
	viper.SetDefault("Identities", map[string]string{})
	viper.AutomaticEnv()

	err := os.MkdirAll(os.ExpandEnv("$HOME/.dapp"), 0744)
	if err != nil {
		errors.Print(err)
		os.Exit(-1)
	}

	err = viper.ReadInConfig()
	if _, ok := err.(*viper.ConfigFileNotFoundError); ok {
		return
	}

	if err != nil {
		errors.Print(err)
		os.Exit(-1)
	}

	err = os.MkdirAll(viper.GetString("CacheDir"), 0744)
	if err != nil {
		errors.Print(err)
		os.Exit(-1)
	}

	id := getIdentity(identity)

	app = dapp.GetApplication(
		"GDCQKPQOB5MSBLHWCNDESVVPQVTRWF2JLCR6LRXTXEMPS3IZCEKY7F6V",
	)

	app.Login(id)
}
