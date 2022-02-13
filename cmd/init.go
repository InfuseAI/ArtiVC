/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
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
	Run: config,
}

func repoInit(cmd *cobra.Command, args []string) {
	cwd, _ := os.Getwd()
	core.InitRepo(cwd)
}

func init() {
	rootCmd.AddCommand(configCommand)
}
