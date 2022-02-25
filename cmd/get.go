package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/infuseai/artiv/internal/core"
	"github.com/spf13/cobra"
)

// getCmd represents the download command
var getCmd = &cobra.Command{
	Use:                   "get [-o <output>] <repository>[@<commit>|<tag>]",
	DisableFlagsInUseLine: true,
	Short:                 "Download data from a repository",
	Example: `  # Download the latest version. The data go to "mydataset" folder.
  art get s3://bucket/mydataset

  # Download the specific version
  art get s3://mybucket/path/to/mydataset@v1.0.0
  
  # Download to a specific folder
  art get -o /tmp/mydataset s3://bucket/mydataset`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		repoUrl, ref, err := parseRepoStr(args[0])
		baseDir, err := cmd.Flags().GetString("output")
		if err != nil {
			exitWithError(err)
			return
		}

		if baseDir == "" {
			comps := strings.Split(repoUrl, "/")
			if len(comps) == 0 {
				exitWithFormat("invlaid path: %v\n", repoUrl)
				return
			}
			baseDir = comps[len(comps)-1]
		}
		baseDir, err = filepath.Abs(baseDir)
		if err != nil {
			exitWithError(err)
			return
		}

		metadataDir, _ := os.MkdirTemp(os.TempDir(), "*-art")
		defer os.RemoveAll(metadataDir)

		config := core.NewConfig(baseDir, metadataDir, repoUrl)

		mngr, err := core.NewArtifactManager(config)
		if err != nil {
			exitWithError(err)
			return
		}

		options := core.PullOptions{NoFetch: true}
		if ref != "" {
			options.RefOrCommit = &ref
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
			options.Mode = core.ChangeModeMerge
		} else if syncMode {
			options.Mode = core.ChangeModeSync
		}

		err = mngr.Pull(options)
		if err != nil {
			exitWithError(err)
			return
		}
	},
}

func init() {
	getCmd.Flags().StringP("output", "o", "", "Output directory")
	getCmd.Flags().Bool("merge", false, "Merge data from the commit. No files would be deleted")
	getCmd.Flags().Bool("sync", false, "Sync data from the commit. The missing files would be deleted")
}
