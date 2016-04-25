package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// buildCmd represents the publish command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "builds and publishes the app to ipfs",
	Long:  `The build command is a lower-level command for building dapps.  It builds a release bundle into a temporary directory then adds that directory to ipfs.`,
	Run: func(cmd *cobra.Command, args []string) {

		dir, err := ioutil.TempDir("", "dapp-build")
		mustSucceed(err)
		defer os.RemoveAll(dir)

		binaryPath := filepath.Join(dir, "bin")

		cargs := append([]string{
			"build", "-o", binaryPath,
		}, args...)

		err = exec.Command("go", cargs...).Run()
		mustSucceed(err)

		stdout, err := exec.Command("ipfs", "add", "-q", binaryPath).Output()
		mustSucceed(err)

		hashes := strings.Split(strings.TrimSpace(string(stdout)), "\n")

		fmt.Println(hashes[len(hashes)-1])
	},
}

func init() {
	RootCmd.AddCommand(buildCmd)
}
