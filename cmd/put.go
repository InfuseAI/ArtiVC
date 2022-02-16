/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/infuseai/art/internal/core"
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "put",
	Short: "Upload files from local to repository",
	Long: `Upload files from local to repository. For example:

# put the latest version
art put ./folder/ /path/to/mydataset
# put the specific version
art put ./folder/ /path/to/mydataset@v1.0.0.`,
	Run:  put,
	Args: cobra.ExactArgs(2),
}

func put(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		log.Fatal("upload require 2 argument")
		os.Exit(1)
	}

	baseDir, err := filepath.Abs(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return
	}

	metadataDir, _ := os.MkdirTemp(os.TempDir(), "*-art")
	repoUrl := args[1]

	config := core.NewConfig(baseDir, metadataDir, repoUrl)

	mngr, err := core.NewArtifactManager(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return
	}

	mngr.Push()
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}
