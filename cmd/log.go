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

var logCommand = &cobra.Command{
	Use:   "log",
	Short: "Log commits",
	Long: `Log commits in the repository. For example:

# list the files for the latest version
art log

# list the files for the specific version
art log v1.0.0`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := core.LoadConfig("")

		var ref string
		if len(args) == 0 {
			ref = core.RefLatest
		} else if len(args) == 1 {
			ref = args[0]
		} else {
			fmt.Fprintf(os.Stderr, "requires 0 or 1 argument\n")
			os.Exit(1)
		}

		if err != nil {
			fmt.Printf("log %v \n", err)
			return
		}

		mngr, err := core.NewArtifactManager(config)
		if err != nil {
			fmt.Printf("log %v \n", err)
			return
		}

		err = mngr.Log(ref)
		if err != nil {
			fmt.Printf("log %v \n", err)
		}
	},
}

func init() {
}
