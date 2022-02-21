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
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull data from the repository to the workspace",
	Long: `Pull data from the repository. For example:

# switch to the latest version
art pull

# switch to the specific version
art pull v1.0.0
`,
	Run: func(cmd *cobra.Command, args []string) {

		config, err := core.LoadConfig("")
		if err != nil {
			fmt.Printf("pull %v \n", err)
			return
		}

		mngr, err := core.NewArtifactManager(config)
		if err != nil {
			fmt.Printf("pull %v \n", err)
			return
		}

		err = mngr.Pull(core.PullOptions{Fetch: true})
		if err != nil {
			fmt.Printf("pull %v \n", err)
		}
	},
}

func init() {
}
