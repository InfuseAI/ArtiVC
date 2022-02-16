/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/infuseai/art/internal/core"
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "put",
	Short: "Upload files from local to repository",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run:  put,
	Args: cobra.ExactArgs(2),
}

func put(cmd *cobra.Command, args []string) {
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

	mngr, err := core.NewArtifactManager(options)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return
	}

	mngr.Push()
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}
