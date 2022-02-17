package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

func exitWithFormat(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: ")
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
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
