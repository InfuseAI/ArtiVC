package repository

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getAbsFilePath(path string) string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	if !strings.HasPrefix(path, "/") {
		path, err = filepath.Abs(filepath.Join(cwd, path))
		if err != nil {
			return ""
		}
	}
	return path
}

func Test(t *testing.T) {
	testCases := []struct {
		desc   string
		repo   string
		scheme string
		host   string
		path   string
		name   string
	}{
		{repo: "/tmp", scheme: "file", host: "", path: getAbsFilePath("/tmp"), name: "tmp"},
		{repo: "tmp", scheme: "file", host: "", path: getAbsFilePath("tmp"), name: "tmp"},
		{repo: "../tmp", scheme: "file", host: "", path: getAbsFilePath("../tmp"), name: "tmp"},
		{repo: "file:///tmp", scheme: "file", host: "", path: "/tmp", name: "tmp"},
		{repo: "host:/tmp", scheme: "ssh", host: "host", path: "/tmp", name: "tmp"},
		{repo: "host:tmp", scheme: "ssh", host: "host", path: "tmp", name: "tmp"},
		{repo: "host:../tmp", scheme: "ssh", host: "host", path: "../tmp", name: "tmp"},
		{repo: "ssh://host/tmp", scheme: "ssh", host: "host", path: "/tmp", name: "tmp"},
		{repo: "xyz://host/tmp", scheme: "xyz", host: "host", path: "/tmp", name: "tmp"},
		{repo: "xyz://host", scheme: "xyz", host: "host", path: "", name: "host"},
	}

	for _, tC := range testCases {
		t.Run("pares repo "+tC.repo, func(t *testing.T) {
			result, err := ParseRepo(tC.repo)
			if err != nil {
				t.Error(err)
				return
			}

			assert.Equal(t, tC.scheme, result.scheme)
			assert.Equal(t, tC.host, result.host)
			assert.Equal(t, tC.path, result.path)

			repoName, err := ParseRepoName(result)
			if err != nil {
				t.Error(err)
				return
			}

			assert.Equal(t, tC.name, repoName)
		})
	}
}
