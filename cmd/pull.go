package cmd

import (
	"errors"

	"github.com/infuseai/artivc/internal/core"
	"github.com/spf13/cobra"
)

// getCmd represents the download command
var pullCmd = &cobra.Command{
	Use:   "pull [<commit>|<tag>]",
	Short: "Pull data from the repository",
	Example: `  # Pull the latest version
  avc pull

  # Pull from a specifc version
  avc pull v1.0.0

  # Pull partial files
  avc pull -- path/to/partia
  avc pull v0.1.0 -- path/to/partia ...`,
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

		option.DryRun, err = cmd.Flags().GetBool("dry-run")
		if err != nil {
			exitWithError(err)
		}

		option.Delete, err = cmd.Flags().GetBool("delete")
		if err != nil {
			exitWithError(err)
		}

		argsLenBeforeDash := cmd.Flags().ArgsLenAtDash()
		if argsLenBeforeDash == -1 {
			if len(args) == 1 {
				option.RefOrCommit = &args[0]
			} else if len(args) > 1 {
				exitWithError(errors.New("please specify \"--\" flag teminator"))
			}
		} else {
			if argsLenBeforeDash == 1 {
				option.RefOrCommit = &args[0]
			}

			if len(args)-argsLenBeforeDash > 0 {
				if option.Delete {
					exitWithError(errors.New("cannot pull partial files and specify delete flag at the same time"))
				}

				fileInclude := core.NewAvcInclude(args[argsLenBeforeDash:])
				option.FileFilter = func(path string) bool {
					return fileInclude.MatchesPath(path)
				}
			}
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
	pullCmd.Flags().Bool("delete", false, "Delete extra files which are not listed in commit")
}
