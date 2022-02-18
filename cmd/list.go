/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/infuseai/art/internal/core"
	"github.com/spf13/cobra"
)

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "List files in the repository",
	Long:  "List files in the repository",
	Example: `  # list the files for the latest version
  art list

  # list the files for the specific version
  art list v1.0.0`,
	Run: func(cmd *cobra.Command, args []string) {
		var ref string
		if len(args) == 0 {
			ref = core.RefLatest
		} else if len(args) == 1 {
			ref = args[0]
		} else {
			fmt.Fprintf(os.Stderr, "requires 0 or 1 argument\n")
			os.Exit(1)
		}

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

		err = mngr.List(ref)
		if err != nil {
			fmt.Printf("list %v \n", err)
		}
	},
}

func init() {
}
