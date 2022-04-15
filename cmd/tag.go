package cmd

import (
	"github.com/infuseai/artivc/internal/core"
	"github.com/spf13/cobra"
)

var tagCommand = &cobra.Command{
	Use:                   "tag [--delete <tag>] [<tag>]",
	DisableFlagsInUseLine: true,
	Short:                 "List or manage tags",
	Example: `  # List the tags
  avc tag

  # Tag the lastest commit
  avc tag v1.0.0

  # Tag the specific commit
  avc tag --ref a1b2c3d4 v1.0.0

  # Delete a tags
  avc tag --delete v1.0.0`,
	Args: cobra.RangeArgs(0, 2),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := core.LoadConfig("")
		exitWithError(err)

		mngr, err := core.NewArtifactManager(config)
		exitWithError(err)

		if len(args) == 0 {
			exitWithError(mngr.ListTags())
		} else if len(args) == 1 {
			tag := args[0]
			refOrCommit, err := cmd.Flags().GetString("ref")
			exitWithError(err)
			delete, err := cmd.Flags().GetBool("delete")
			exitWithError(err)

			if !delete {
				exitWithError(mngr.AddTag(refOrCommit, tag))
			} else {
				exitWithError(mngr.DeleteTag(tag))
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
