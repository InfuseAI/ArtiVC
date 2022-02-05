/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/infuseai/art/internal"
	"github.com/spf13/cobra"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		Upload()
	},
}

func Upload() {
	fileList := make([]string, 0)

	filepath.Walk("/Users/popcorny/art/myart", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		fileList = append(fileList, path)
		return nil
	})

	for _, path := range fileList {
		fmt.Printf("%s\t%s\n", internal.Sha1sum(path), path)
	}
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}
