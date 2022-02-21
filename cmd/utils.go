package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

const (
	GROUP_BASIC     = "basic"
	GROUP_WORKSPACE = "workspace"
)

func exitWithError(err error) {
	cobra.CheckErr(err)
}

func exitWithFormat(format string, a ...interface{}) {
	cobra.CheckErr(fmt.Sprintf(format, a...))
}

func parseRepoStr(repoAndRef string) (repoUrl string, ref string, err error) {
	comps := strings.Split(repoAndRef, "@")
	if len(comps) == 1 {
		repoUrl = repoAndRef
	} else if len(comps) == 2 {
		repoUrl = comps[0]
		ref = comps[1]
	} else {
		err = errors.New("Invalid repository: " + repoAndRef)
	}
	return
}
