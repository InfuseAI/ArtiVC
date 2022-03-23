package cmd

import (
	"github.com/infuseai/artivc/internal/core"
	"github.com/spf13/cobra"
)

var logCommand = &cobra.Command{
	Use:                   "log [<commit>|<tag>]",
	DisableFlagsInUseLine: true,
	Short:                 "Log commits",
	Example: `  # Log commits from the latest
  avc log

  # Log commits from a specific version
  avc log v1.0.0`,
	Args: cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := core.LoadConfig("")

		var ref string
		if len(args) == 0 {
			ref = core.RefLatest
		} else {
			ref = args[0]
		}

		if err != nil {
			exitWithError(err)
			return
		}

		mngr, err := core.NewArtifactManager(config)
		if err != nil {
			exitWithError(err)
			return
		}

		err = mngr.Log(ref)
		if err != nil {
			exitWithError(err)
			return
		}
	},
}

func init() {
}
