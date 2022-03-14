package cmd

import (
	"errors"
	"fmt"
	"io"
	neturl "net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const (
	GROUP_BASIC = "basic"
	GROUP_QUICK = "quick"
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

func parseRepoName(repoUrl string) (string, error) {
	url, err := neturl.Parse(repoUrl)
	if err != nil {
		return "", err
	}

	if url.Path == "" {
		return url.Hostname(), nil
	}

	name := filepath.Base(url.Path)
	if name == "/" {
		return url.Hostname(), nil
	}

	return name, nil
}

func transformRepoUrl(base string, repo string) (string, error) {
	url, err := neturl.Parse(repo)
	if err != nil {
		return "", err
	}

	if url.Scheme != "" {
		return repo, nil
	}

	if strings.HasPrefix(repo, "/") {
		return repo, nil
	}

	return filepath.Abs(filepath.Join(base, url.Path))
}

func isDirEmpty(dir string) bool {
	f, err := os.Open(dir)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true
	}
	return false
}
