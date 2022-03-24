package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docsCommand = &cobra.Command{
	Use:   "docs",
	Short: "Generate docs",
	Long: `Generate docs. For example:

avc docs`,
	Run: func(cmd *cobra.Command, args []string) {
		const DocDir = "./generated_docs"
		err := os.Mkdir(DocDir, fs.ModePerm)

		if err == nil || (err != nil && os.IsExist(err)) {
			// pass when directory existing
		} else {
			fmt.Printf("Failed to create %s, skip to generate documents\n", DocDir)
			return
		}
		linkHandler := func(name string) string {
			base := strings.TrimSuffix(name, path.Ext(name))
			return "/commands/" + strings.ToLower(base) + "/"
		}

		doc.GenMarkdownTreeCustom(cmd.Root(), DocDir, func(filestring string) string { return "" }, linkHandler)
	},
}
