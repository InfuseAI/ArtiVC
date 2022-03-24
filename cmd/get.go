package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/infuseai/artivc/internal/core"
	"github.com/spf13/cobra"
)

// getCmd represents the download command
var getCmd = &cobra.Command{
	Use:                   "get [-o <output>] <repository>[@<commit>|<tag>] [--] <pathspec>...",
	DisableFlagsInUseLine: true,
	Short:                 "Download data from a repository",
	Example: `  # Download the latest version. The data go to "mydataset" folder.
  avc get s3://bucket/mydataset

  # Download the specific version
  avc get s3://mybucket/path/to/mydataset@v1.0.0
  
  # Download to a specific folder
  avc get -o /tmp/mydataset s3://bucket/mydataset

  # Download partial files
  avc get -o /tmp/mydataset s3://bucket/mydataset -- path/to/file1 path/to/file2 data/`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		repoUrl, ref, err := parseRepoStr(args[0])
		if err != nil {
			exitWithError(err)
			return
		}
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

		metadataDir, _ := os.MkdirTemp(os.TempDir(), "*-avc")
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

		options.Delete, err = cmd.Flags().GetBool("delete")
		if err != nil {
			exitWithError(err)
		}

		if len(args) > 1 {
			if options.Delete {
				exitWithError(errors.New("cannot download partial files and specify delete flag at the same time"))
			}
			fileInclude := core.NewAvcInclude(args[1:])
			options.FileFilter = func(path string) bool {
				return fileInclude.MatchesPath(path)
			}
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
	getCmd.Flags().Bool("delete", false, "Delete extra files which are not listed in commit")
}
