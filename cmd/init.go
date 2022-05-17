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
		cwd, err := os.Getwd()
		exitWithError(err)

		result, err := repository.ParseRepo(args[0])
		exitWithError(err)
		repo := result.Repo

		if strings.HasPrefix(repo, "http") && !repository.IsAzureStorageUrl(repo) {
			exitWithError(errors.New("init not support under http(s) repo"))
		}

		_, err = repository.NewRepository(result)
		exitWithError(err)

		fmt.Printf("Initialize the artivc workspace of the repository '%s'\n", repo)
		exitWithError(core.InitWorkspace(cwd, repo))
	},
}

func init() {
}
