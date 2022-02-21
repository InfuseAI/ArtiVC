package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

var (
	version = "v0.1-dev"

	// overwrite version when tagVersion exists
	tagVersion = ""

	// gitCommit is the git sha1
	gitCommit = ""

	// gitTreeState is the state of the git tree {dirty or clean}
	gitTreeState = ""
)

type BuildInfo struct {
	Version      string
	GitCommit    string
	GitTreeState string
	GoVersion    string
}

func GetVersion() string {
	info := BuildInfo{
		Version:      version,
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		GoVersion:    runtime.Version(),
	}

	if tagVersion != "" {
		info.Version = tagVersion
	}

	data, _ := json.Marshal(info)
	return fmt.Sprintf("version.BuildInfo%s", string(data))
}

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long: `Print the version information. For example:

art version
version.BuildInfo{"Version":"v0.1-dev","GitCommit":"59b5c650fbed4d91c1e54b7cb3c3f6f0c50e5fa4","GitTreeState":"dirty","GoVersion":"go1.17.5"}
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(GetVersion())
	},
}

func init() {

}
