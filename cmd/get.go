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

// getCmd represents the download command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: get,
}

func get(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		log.Fatal("upload require 2 argument")
		os.Exit(1)
	}

	src := args[0]
	dest := args[1]

	options := core.ArtifactManagerOptions{
		BaseDir:    &dest,
		Repository: &src,
	}

	mngr, err := core.NewArtifactManager(options)
	if err != nil {
		fmt.Printf("pull %v \n", err)
		return
	}

	err = mngr.Pull()
	if err != nil {
		fmt.Printf("pull %v \n", err)
		return
	}
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
