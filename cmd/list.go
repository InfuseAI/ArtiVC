package cmd

import (
	"github.com/infuseai/artiv/internal/core"
	"github.com/spf13/cobra"
)

var listCommand = &cobra.Command{
	Use:     "list",
	Short:   "List files of a commit",
	Aliases: []string{"ls"},
	Example: `  # List files for the latest version
  art list

  # List files for the specific version
  art list v1.0.0`,
	Args: cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		var ref string
		if len(args) == 0 {
			ref = core.RefLatest
		} else {
			ref = args[0]
		}

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

		err = mngr.List(ref)
		if err != nil {
			exitWithError(err)
		}
	},
}

func init() {
}
