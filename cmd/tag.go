/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/infuseai/art/internal/core"
	"github.com/spf13/cobra"
)

var tagCommand = &cobra.Command{
	Use:   "tag",
	Short: "Manage tags in the repository",
	Long: `Manage tags in the repository. For example:

# list the tags
art tag

# Add a tag
art tag v1.0.0

# Delete a tags
art tag --delete v1.0.0
`,
	Run: tag,
}

func tag(cmd *cobra.Command, args []string) {
	config, err := core.LoadConfig("")
	if err != nil {
		exitWithError(err)
	}

	mngr, err := core.NewArtifactManager(config)
	if err != nil {
		fmt.Printf("log %v \n", err)
		return
	}

	if len(args) == 0 {
		err := mngr.ListTags()
		if err != nil {
			exitWithError(err)
		}
	} else if len(args) == 1 {
		tag := args[0]
		refOrCommit, err := cmd.Flags().GetString("ref")
		if err != nil {
			exitWithError(err)
		}
		delete, err := cmd.Flags().GetBool("delete")
		if err != nil {
			exitWithError(err)
		}

		if !delete {
			err := mngr.AddTag(refOrCommit, tag)
			if err != nil {
				exitWithError(err)
			}
		} else {
			err := mngr.DeleteTag(tag)
			if err != nil {
				exitWithError(err)
			}
		}
	} else {
		exitWithFormat("requires 0 or 1 argument\n")
	}
}

func init() {
	rootCmd.AddCommand(tagCommand)

	tagCommand.Flags().BoolP("delete", "D", false, "Delete a tag")
	tagCommand.Flags().String("ref", core.RefLatest, "The source commit or reference to be tagged")
}
