package cmd

import (
	"errors"
	"os"
	"strings"

	"github.com/infuseai/artiv/internal/core"
	"github.com/infuseai/artiv/internal/repository"
	"github.com/spf13/cobra"
)

var initCommand = &cobra.Command{
	Use:                   "init <repository>",
	Short:                 "Initiate a workspace",
	DisableFlagsInUseLine: true,
	Example: `  # Init a workspace with local repository
  art init /path/to/mydataset

  # Init a workspace with s3 repoisotry
  art init s3://mybucket/path/to/mydataset`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cwd, _ := os.Getwd()
		repo := args[0]

		if strings.HasPrefix(repo, "http") {
			exitWithError(errors.New("init not support under http(s) repo"))
			return
		}

		_, err := repository.NewRepository(repo)
		if err != nil {
			exitWithError(err)
			return
		}

		core.InitWorkspace(cwd, repo)
	},
}

func init() {
}
