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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: repoInit,
}

func repoInit(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("init requires 1 argument")
		os.Exit(1)
	}

	cwd, _ := os.Getwd()
	repo := args[0]
	core.InitRepo(cwd, repo)
}

func init() {
	rootCmd.AddCommand(initCommand)
}
