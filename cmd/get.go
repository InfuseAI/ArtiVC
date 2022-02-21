/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/infuseai/art/internal/core"
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

		if len(args) != 1 {
			log.Fatal("get require only 1 argument")
			os.Exit(1)
		}

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

		options := core.PullOptions{}
		if ref != "" {
			options.RefOrCommit = &ref
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
}
