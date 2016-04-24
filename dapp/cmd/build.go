package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// buildCmd represents the publish command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "builds and publishes the app to ipfs",
	Long:  `The build command is a lower-level command for building dapps.  It builds a release bundle into a temporary directory then adds that directory to ipfs.`,
	Run: func(cmd *cobra.Command, args []string) {

		dir, err := ioutil.TempDir("", "dapp-build")
		if err != nil {
			errors.Print(err)
			os.Exit(-1)
		}
		defer os.RemoveAll(dir)

		binaryPath := filepath.Join(dir, "bin")

		cargs := append([]string{
			"build", "-o", binaryPath,
		}, args...)

		err = exec.Command("go", cargs...).Run()
		if err != nil {
			errors.Print(err)
			os.Exit(-1)
		}

		stdout, err := exec.Command("ipfs", "add", "-q", binaryPath).Output()
		if err != nil {
			errors.Print(err)
			os.Exit(-1)
		}

		hashes := strings.Split(strings.TrimSpace(string(stdout)), "\n")

		fmt.Println(hashes[len(hashes)-1])
	},
}

func init() {
	RootCmd.AddCommand(buildCmd)
}
