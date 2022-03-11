package cmd

import (
	"errors"
	"os"
	"strings"

	"github.com/infuseai/artiv/internal/core"
	"github.com/infuseai/artiv/internal/repository"
	"github.com/spf13/cobra"
)

var cloneCommand = &cobra.Command{
	Use:                   "clone <repository>",
	Short:                 "clone a workspace",
	DisableFlagsInUseLine: true,
	Example: `  # clone a workspace with local repository
  art clone /path/to/mydataset

  # clone a workspace with s3 repository
  art clone s3://mybucket/path/to/mydataset`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cwd, _ := os.Getwd()
		repo, err := transformRepoUrl(cwd, args[0])
		if err != nil {
			exitWithError(err)
			return
		}

		if strings.HasPrefix(repo, "http") {
			exitWithError(errors.New("clone not support under http(s) repo"))
			return
		}

		_, err = repository.NewRepository(repo)
		if err != nil {
			exitWithError(err)
			return
		}

		core.InitWorkspace(cwd, repo)

		config, err := core.LoadConfig("")
		if err != nil {
			exitWithError(err)
			return
		}

		mngr, err := core.NewArtifactManager(config)
		if err != nil {
			exitWithError(err)
			return
		}

		option := core.PullOptions{}
		err = mngr.Pull(option)
		if err != nil {
			exitWithError(err)
			return
		}
	},
}

func init() {
}
