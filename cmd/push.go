/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/infuseai/art/internal/core"
	"github.com/spf13/cobra"
)

// getCmd represents the download command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Make all changes to a commit in the repository",
	Long: `Make all changes to a commit in the repository. For example:

# push current folder to remote
art push -m 'This is initial version'`,
	Run: push,
}

func push(cmd *cobra.Command, args []string) {

	config, err := core.LoadConfig("")
	if err != nil {
		fmt.Printf("pull %v \n", err)
		return
	}

	mngr, err := core.NewArtifactManager(config)
	if err != nil {
		fmt.Printf("push %v \n", err)
		return
	}

	err = mngr.Push()
	if err != nil {
		fmt.Printf("push %v \n", err)
	}
}

func init() {
	rootCmd.AddCommand(pushCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
