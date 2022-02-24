package cmd

import (
	"os"

	"github.com/infuseai/art/internal/core"
	"github.com/infuseai/art/internal/repository"
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
