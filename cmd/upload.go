/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"
	"os"

	"github.com/infuseai/art/internal/core"
	"github.com/spf13/cobra"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload files from local to repository",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run:  Upload,
	Args: cobra.ExactArgs(2),
}

func Upload(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		log.Fatal("upload require 2 argument")
		os.Exit(1)
	}

	src := args[0]
	dest := args[1]

	options := core.ArtifactManagerOptions{
		BaseDir:    &src,
		Repository: &dest,
	}

	mngr := core.NewArtifactManager(options)
	mngr.Push()
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}
