package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sshParseListRecord(t *testing.T) {

}

func Test(t *testing.T) {
	testCases := []struct {
		desc  string
		in    string
		name  string
		isDir bool
	}{
		{
			desc:  "directory",
			in:    "drwxrwxrwt 1 root root 4096 Mar 25 02:45 /tmp",
			name:  "/tmp",
			isDir: true,
		},
		{
			desc:  "file",
			in:    "-rw-r--r-- 1 root root   29 Mar 25 01:24 a",
			name:  "a",
			isDir: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			info, err := sshParseListRecord(tC.in)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tC.name, info.Name())
			assert.Equal(t, tC.isDir, info.IsDir())
		})
	}
}
