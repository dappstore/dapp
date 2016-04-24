package cmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "dapp",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
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

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dapp.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/dapp")
	viper.AddConfigPath("$HOME/.dapp")
	viper.SetDefault("CacheDir", os.ExpandEnv("$HOME/.dapp/cache"))
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
}
