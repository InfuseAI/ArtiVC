/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/infuseai/art/internal/core"
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

		option := core.PullOptions{Fetch: true}
		option.DryRun, err = cmd.Flags().GetBool("dry-run")
		if err != nil {
			exitWithError(err)
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
}
