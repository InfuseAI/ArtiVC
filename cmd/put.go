/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/infuseai/art/internal/core"
	"github.com/spf13/cobra"
)

var putCmd = &cobra.Command{
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
	// arguments
	if len(args) != 2 {
		log.Fatal("upload require 2 argument")
		os.Exit(1)
	}

	baseDir, err := filepath.Abs(args[0])
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

	// Create temp metadata
	metadataDir, _ := os.MkdirTemp(os.TempDir(), "*-art")
	defer os.RemoveAll(metadataDir)

	repoUrl := args[1]

	config := core.NewConfig(baseDir, metadataDir, repoUrl)

	// push
	mngr, err := core.NewArtifactManager(config)
	if err != nil {
		exitWithError(err)
	}

	mngr.Push(option)
}

func init() {
	rootCmd.AddCommand(putCmd)

	putCmd.Flags().StringP("message", "m", "", "Commit meessage")
}
