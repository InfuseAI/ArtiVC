package cmd

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docsCommand = &cobra.Command{
	Use:   "docs",
	Short: "Generate docs",
	Long: `Generate docs. For example:

art docs`,
	Run: func(cmd *cobra.Command, args []string) {
		const DocDir = "./generated_docs"
		err := os.Mkdir(DocDir, fs.ModePerm)

		if err == nil || (err != nil && os.IsExist(err)) {
			// pass when directory existing
		} else {
			fmt.Printf("Failed to create %s, skip to generate documents\n", DocDir)
			return
		}
		doc.GenMarkdownTree(cmd.Root(), DocDir)
	},
}
