package cmd

import (
	"github.com/infuseai/artiv/internal/core"
	"github.com/spf13/cobra"
)

// getCmd represents the download command
var pullCmd = &cobra.Command{
	Use:   "pull [<commit>|<tag>]",
	Short: "Pull data from the repository",
	Example: `  # Pull the latest version
  art pull

  # Pull from a specifc version
  art pull v1.0.0`,
	Args: cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
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

		// options
		option := core.PullOptions{}
		if len(args) > 0 {
			option.RefOrCommit = &args[0]
		}

		option.DryRun, err = cmd.Flags().GetBool("dry-run")
		if err != nil {
			exitWithError(err)
		}

		mergeMode, err := cmd.Flags().GetBool("merge")
		if err != nil {
			exitWithError(err)
		}

		syncMode, err := cmd.Flags().GetBool("sync")
		if err != nil {
			exitWithError(err)
		}

		if mergeMode && syncMode {
			exitWithFormat("only one of --merge and --sync can be set")
		}
		if mergeMode {
			option.Mode = core.ChangeModeMerge
		} else if syncMode {
			option.Mode = core.ChangeModeSync
		}

		err = mngr.Pull(option)
		if err != nil {
			exitWithError(err)
			return
		}
	},
}

func init() {
	pullCmd.Flags().Bool("dry-run", false, "Dry run")
	pullCmd.Flags().Bool("merge", false, "Merge data from the commit. No files would be deleted")
	pullCmd.Flags().Bool("sync", false, "Sync data from the commit. The missing files would be deleted")
}
