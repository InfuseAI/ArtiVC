/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/infuseai/art/internal/core"
	"github.com/spf13/cobra"
)

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "List files in the repository",
	Long: `List files in the repository. For example:

# list the files for the latest version
art list

# list the files for the specific version
art list v1.0.0`,
	Run: list,
}

func list(cmd *cobra.Command, args []string) {
	config, err := core.LoadConfig("")
	if err != nil {
		fmt.Printf("list %v \n", err)
		return
	}

	mngr, err := core.NewArtifactManager(config)
	if err != nil {
		fmt.Printf("list %v \n", err)
		return
	}

	err = mngr.List()
	if err != nil {
		fmt.Printf("list %v \n", err)
	}
}

func init() {
	rootCmd.AddCommand(listCommand)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
