package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	testCases := []struct {
		desc   string
		repo   string
		scheme string
		host   string
		path   string
	}{
		{repo: "/tmp", scheme: "file", host: "", path: "/tmp"},
		{repo: "tmp", scheme: "file", host: "", path: "tmp"},
		{repo: "../tmp", scheme: "file", host: "", path: "../tmp"},
		{repo: "file:///tmp", scheme: "file", host: "", path: "/tmp"},
		{repo: "host:/tmp", scheme: "ssh", host: "host", path: "/tmp"},
		{repo: "host:tmp", scheme: "ssh", host: "host", path: "tmp"},
		{repo: "host:../tmp", scheme: "ssh", host: "host", path: "../tmp"},
		{repo: "ssh://host/tmp", scheme: "ssh", host: "host", path: "/tmp"},
		{repo: "xyz://host/tmp", scheme: "xyz", host: "host", path: "/tmp"},
	}

	for _, tC := range testCases {
		t.Run("pares repo "+tC.repo, func(t *testing.T) {
			result, err := parseRepo(tC.repo)
			if err != nil {
				t.Error(err)
				return
			}

			assert.Equal(t, tC.scheme, result.scheme)
			assert.Equal(t, tC.host, result.host)
			assert.Equal(t, tC.path, result.path)
		})
	}
}
