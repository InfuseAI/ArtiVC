/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/infuseai/art/internal/core"
	"github.com/infuseai/art/internal/repository"
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
	Run:  Upload,
	Args: cobra.ExactArgs(2),
}

func Upload(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		log.Fatal("upload require 2 argument")
		os.Exit(1)
	}

	src := args[0]
	dest := args[1]

	commit := core.Commit{
		CreatedAt: time.Now(),
		Message:   nil,
		Blobs:     make([]core.BlobMetaData, 0),
	}

	filepath.Walk(src, func(absPath string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", absPath, err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		path := absPath[len(src)+1:]
		metadata, err := core.MakeBlobMetadata(src, path)
		if err != nil {
			log.Fatalf("cannot make metadata: %s", path)
			return err
		}

		commit.Blobs = append(commit.Blobs, metadata)
		return nil
	})

	repo := repository.LocalFileSystemRepository{
		BaseDir: src,
		RepoDir: dest,
	}

	for _, metadata := range commit.Blobs {
		log.Printf("upload %s\n", metadata.Path)
		err := repo.UploadBlob(metadata)
		if err != nil {
			log.Fatalf("cannot upload blob: %s\n", metadata.Path)
			break
		}
	}

	_, hash := core.MakeCommitMetadata(&commit)
	repo.Commit(commit)
	repo.AddRef("latest", hash)
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}
