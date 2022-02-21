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

var configCommand = &cobra.Command{
	Use:                   "config [<key> [<value>]]",
	Short:                 "Configure the workspace",
	Long:                  "Configure the workspace. The config file is stored at \".art/config\".",
	DisableFlagsInUseLine: true,
	Example: `  # List the config
  art config

  # Get the config
  art config repo.url

  # Set the config
  art config repo.url s3://your-bucket/data`,
	Args: cobra.RangeArgs(0, 2),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := core.LoadConfig("")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}

		switch len(args) {
		case 0:
			config.Print()
		case 1:
			value := config.Get(args[0])
			if value != nil {
				fmt.Println(value)
			} else {
				fmt.Fprintf(os.Stderr, "key not found: %s\n", args[0])
				os.Exit(1)
			}
		case 2:
			config.Set(args[0], args[1])
			err := config.Save()
			if err != nil {
				fmt.Println(err)
			}
		}
	},
}

func init() {
}
