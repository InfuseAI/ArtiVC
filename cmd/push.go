package cmd

import (
	"github.com/infuseai/artivc/internal/core"
	"github.com/spf13/cobra"
)

// getCmd represents the download command
var pushCmd = &cobra.Command{
	Use:                   "push [-m <message>]",
	DisableFlagsInUseLine: true,
	Short:                 "Push data to the repository",
	Long:                  `Push data to the repository. There is no branch implemented yet, all put and push commands are always creating a commit and treat as the latest commit.`,
	Example: `  # Push to the latest version
  avc push -m 'Initial version'

  # Push to the latest version and tag to specific version
  avc push -m 'Initial version'
  avc tag v1.0.0`,
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

		option.DryRun, err = cmd.Flags().GetBool("dry-run")
		if err != nil {
			exitWithError(err)
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
	pushCmd.Flags().Bool("dry-run", false, "Dry run")
}
