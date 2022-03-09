package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	baseDir := "/tmp/artiv"
	testCases := []struct {
		desc string
		in   string
		out  string
	}{
		{desc: "local file", in: "/this/is/my/path", out: "/this/is/my/path"},
		{desc: "relative path", in: "../path", out: "/tmp/path"},
		{desc: "relative path2", in: "../../../path", out: "/path"},
		{desc: "normal url (file)", in: "file://mybucket/this/is/my/path", out: "file://mybucket/this/is/my/path"},
		{desc: "normal url (s3)", in: "s3://mybucket/this/is/my/path", out: "s3://mybucket/this/is/my/path"},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			result, err := transformRepoUrl(baseDir, tC.in)
			if err != nil {
				assert.Empty(t, tC.out)
			} else {
				assert.Equal(t, tC.out, result)
			}
		})
	}
}
