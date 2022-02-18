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

var initCommand = &cobra.Command{
	Use:   "init",
	Short: "Initialize a workspace",
	Long: `Initialize a workspace. For example:

cd mydataset/
art init s3://mybucket/path/to/mydataset`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("init requires 1 argument")
			os.Exit(1)
		}

		cwd, _ := os.Getwd()
		repo := args[0]
		core.InitWorkspace(cwd, repo)
	},
}

func init() {
}
