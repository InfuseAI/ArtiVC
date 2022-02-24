package cmd

import (
	"github.com/infuseai/artiv/internal/core"
	"github.com/spf13/cobra"
)

var diffCommand = &cobra.Command{
	Use:   "diff",
	Short: "Diff workspace/commits/references",
	Example: `# Diff two version
art diff v0.1.0 v0.2.0`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		left := args[0]
		right := args[1]
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

		_, err = mngr.Diff(core.DiffOptions{
			LeftRef:  left,
			RightRef: right,
		})
		if err != nil {
			exitWithError(err)
		}
	},
}

func init() {
}
