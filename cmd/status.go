package cmd

import (
	"github.com/infuseai/artiv/internal/core"
	"github.com/spf13/cobra"
)

var statusCommand = &cobra.Command{
	Use:                   "status",
	Short:                 "Diff between the remote repository and the workspace",
	DisableFlagsInUseLine: true,
	Example: `	# check current differences
	art status`,
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

		result, err := mngr.Status()
		if err != nil {
			exitWithError(err)
		}

		result.Print(true)
	},
}

func init() {
}
