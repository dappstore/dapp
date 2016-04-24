package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [ipfs hash]",
	Short: "run downloads the binary specified and runs it",
	Long:  `The run command is a helper command that makes it simple to run any dapp-enabled application.  It downloads an ipfs path into a temporary directory then runs it in the current directory`,
	Run: func(cmd *cobra.Command, args []string) {

		hash := args[0]
		cargs := args[1:]
		bin := filepath.Join(viper.GetString("CacheDir"), hash)
		out, err := exec.Command("ipfs", "get", "-o", bin, hash).CombinedOutput()
		if err != nil {
			fmt.Print(string(out))
			os.Exit(-1)
		}

		err = os.Chmod(bin, 0755)
		if err != nil {
			errors.Print(err)
			os.Exit(-1)
		}

		err = syscall.Exec(bin, cargs, os.Environ())
		if err != nil {
			errors.Print(err)
			os.Exit(-1)
		}
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
}
