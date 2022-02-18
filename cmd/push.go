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
	Use:   "push",
	Short: "Make all changes to a commit in the repository",
	Long: `Make all changes to a commit in the repository. For example:

# push current folder to remote
art push -m 'This is initial version'`,
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
