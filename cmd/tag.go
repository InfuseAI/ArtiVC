package cmd

import (
	"fmt"

	"github.com/infuseai/art/internal/core"
	"github.com/spf13/cobra"
)

var tagCommand = &cobra.Command{
	Use:                   "tag [--delete <tag>] [<tag>]",
	DisableFlagsInUseLine: true,
	Short:                 "List or manage tags",
	Example: `  # List the tags
  art tag

  # Tag the lastest commit
  art tag v1.0.0

  # Tag the specific commit
  art tag --ref a1b2c3d4 v1.0.0  

  # Delete a tags
  art tag --delete v1.0.0`,
	Args: cobra.RangeArgs(0, 2),
	Run: func(cmd *cobra.Command, args []string) {
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
	},
}

func init() {
	tagCommand.Flags().BoolP("delete", "D", false, "Delete a tag")
	tagCommand.Flags().String("ref", core.RefLatest, "The source commit or reference to be tagged")
}
