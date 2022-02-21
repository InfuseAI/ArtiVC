/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/infuseai/art/internal/core"
	"github.com/spf13/cobra"
)

// getCmd represents the download command
var pushCmd = &cobra.Command{
	Use:                   "push [-m <message>]",
	DisableFlagsInUseLine: true,
	Short:                 "Push data to the repository",
	Long:                  `Push data to the repository. There is no branch implemented yet, all put and push commands are always creating a commit and treat as the latest commit.`,
	Example: `  # Push to the latest version
  art push -m 'Initial version'

  # Push to the latest version and tag to specific version
  art push -m 'Initial version'
  art tag v1.0.0`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := core.LoadConfig("")
		if err != nil {
			exitWithError(err)
		}

		// options
		option := core.PushOptions{}
		message, err := cmd.Flags().GetString("message")
		if err != nil {
			exitWithError(err)
		}
		if message != "" {
			option.Message = &message
		}

		// push
		mngr, err := core.NewArtifactManager(config)
		if err != nil {
			exitWithError(err)
		}

		err = mngr.Push(option)
		if err != nil {
			exitWithError(err)
		}
	},
}

func init() {
	pushCmd.Flags().StringP("message", "m", "", "Commit meessage")
}
