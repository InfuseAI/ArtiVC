package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/infuseai/artivc/internal/core"
	"github.com/infuseai/artivc/internal/repository"
	"github.com/spf13/cobra"
)

var initCommand = &cobra.Command{
	Use:                   "init <repository>",
	Short:                 "Initiate a workspace",
	DisableFlagsInUseLine: true,
	Example: `  # Init a workspace with local repository
  avc init /path/to/mydataset

  # Init a workspace with s3 repository
  avc init s3://mybucket/path/to/mydataset`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cwd, _ := os.Getwd()
		repo, err := transformRepoUrl(cwd, args[0])
		if err != nil {
			exitWithError(err)
			return
		}

		if strings.HasPrefix(repo, "http") {
			exitWithError(errors.New("init not support under http(s) repo"))
			return
		}

		_, err = repository.NewRepository(repo)
		if err != nil {
			exitWithError(err)
			return
		}

		fmt.Printf("Initialize the artivc workspace of the repository '%s'\n", repo)
		core.InitWorkspace(cwd, repo)
	},
}

func init() {
}
