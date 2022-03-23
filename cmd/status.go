package cmd

import (
	"fmt"

	"github.com/infuseai/artivc/internal/core"
	"github.com/spf13/cobra"
)

var statusCommand = &cobra.Command{
	Use:                   "status",
	Short:                 "Show the status of the workspace",
	DisableFlagsInUseLine: true,
	Example: `	# check current status
	avc status`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := core.LoadConfig("")
		if err != nil {
			exitWithError(err)
		}

		mngr, err := core.NewArtifactManager(config)
		if err != nil {
			exitWithError(err)
		}

		err = mngr.Fetch()
		if err != nil {
			exitWithError(err)
		}

		fmt.Printf("workspace of the repository '%s'\n\n", config.RepoUrl())

		result, err := mngr.Status()
		if err != nil {
			exitWithError(err)
		}

		result.Print(true)
	},
}

func init() {
}
