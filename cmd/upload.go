/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"

	"github.com/infuseai/art/internal"
	"github.com/spf13/cobra"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload files from local to repository",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: Upload,
}

func Upload(cmd *cobra.Command, args []string) {
	src := "/Users/popcorny/art/myart"
	dest := "/Users/popcorny/art/myrepo"

	fileList := make([]string, 0)

	filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
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
		fmt.Printf("%s\t%s\n", internal.Sha1SumFromFile(path), path)
	}
}

func upload_file(srcFile string, repo string) {
	input, err := ioutil.ReadFile(srcFile)
	hashsum := internal.Sha1SumFromFile(srcFile)

	destFile := fmt.Sprintf(rep)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile(destFile, input, 0644)
	if err != nil {
		fmt.Println("Error creating", destFile)
		fmt.Println(err)
		return
	}
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}
