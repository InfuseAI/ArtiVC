/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/infuseai/art/internal/core"
	"github.com/spf13/cobra"
)

// getCmd represents the download command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Download data from repository",
	Long: `Download data from repository. For example:

# download to 'mydataset' folder
art get /path/to/mydataset
art get file:///path/to/mydataset
art get s3://mybucket/path/to/mydataset`,
	Run: get,
}

func get(cmd *cobra.Command, args []string) {
	var err error

	if len(args) != 1 {
		log.Fatal("get require only 1 argument")
		os.Exit(1)
	}

	repoUrl := args[0]
	baseDir, err := cmd.Flags().GetString("output")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: i%v\n", err)
		return
	}

	if baseDir == "" {
		comps := strings.Split(repoUrl, "/")
		if len(comps) == 0 {
			fmt.Fprintf(os.Stderr, "error: invlaid path: %v\n", repoUrl)
			return
		}
		baseDir = comps[len(comps)-1]
	}
	baseDir, err = filepath.Abs(baseDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return
	}

	metadataDir, _ := os.MkdirTemp(os.TempDir(), "*-art")
	defer os.RemoveAll(metadataDir)

	config := core.NewConfig(baseDir, metadataDir, repoUrl)

	mngr, err := core.NewArtifactManager(config)
	if err != nil {
		fmt.Printf("pull %v \n", err)
		return
	}

	err = mngr.Pull(core.PullOptions{})
	if err != nil {
		fmt.Printf("pull %v \n", err)
		return
	}
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringP("output", "o", "", "Output directory")
}
