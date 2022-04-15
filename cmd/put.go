package cmd

import (
	"os"
	"path/filepath"

	"github.com/infuseai/artivc/internal/core"
	"github.com/spf13/cobra"
)

var putCmd = &cobra.Command{
	Use:                   "put [-m <message>] <dir> <repository>[@<tag>]",
	DisableFlagsInUseLine: true,
	Short:                 "Upload data to a repository",
	Example: `  # Upload the latest version
  avc put ./folder/ /path/to/mydataset

  # Upload the specific version
  avc put ./folder/ /path/to/mydataset@v1.0.0`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		baseDir, err := filepath.Abs(args[0])
		exitWithError(err)

		repoUrl, ref, err := parseRepoStr(args[1])
		exitWithError(err)

		// options
		option := core.PushOptions{}
		message, err := cmd.Flags().GetString("message")
		exitWithError(err)

		if message != "" {
			option.Message = &message
		}
		if ref != "" {
			option.Tag = &ref
		}

		// Create temp metadata
		metadataDir, _ := os.MkdirTemp(os.TempDir(), "*-avc")
		defer os.RemoveAll(metadataDir)

		config := core.NewConfig(baseDir, metadataDir, repoUrl)

		// push
		mngr, err := core.NewArtifactManager(config)
		exitWithError(err)

		exitWithError(mngr.Push(option))
	},
}

func init() {
	putCmd.Flags().StringP("message", "m", "", "Commit meessage")
}
