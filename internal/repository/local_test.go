package repository

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalUpload(t *testing.T) {
	testCases := []struct {
		desc string
		data string
	}{
		{
			desc: "empty file", data: "",
		},
		{
			desc: "non empty file", data: "hello",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			repoDir := t.TempDir()
			tmpDir := t.TempDir()

			repo, err := NewLocalFileSystemRepository(repoDir)
			if err != nil {
				t.Error(err)
			}

			err = os.WriteFile(tmpDir+"/test", []byte(tC.data), 0644)
			if err != nil {
				t.Error(err)
			}

			err = repo.Upload(tmpDir+"/test", "path/to/the/test", nil)
			if err != nil {
				t.Error(err)
			}
			data, err := os.ReadFile(repoDir + "/path/to/the/test")
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, []byte(tC.data), []byte(data))
		})
	}
}

func TestLocalDownload(t *testing.T) {
	testCases := []struct {
		desc string
		data string
	}{
		{
			desc: "empty file", data: "",
		},
		{
			desc: "non empty file", data: "hello",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			repoDir := t.TempDir()
			tmpDir := t.TempDir()

			repo, err := NewLocalFileSystemRepository(repoDir)
			if err != nil {
				t.Error(err)
			}

			err = os.MkdirAll(repoDir+"/path/to/the", os.ModePerm)
			if err != nil {
				t.Error(err)
			}

			err = os.WriteFile(repoDir+"/path/to/the/test", []byte(tC.data), 0644)
			if err != nil {
				t.Error(err)
			}

			err = repo.Download("path/to/the/test", tmpDir+"/test", nil)
			if err != nil {
				t.Error(err)
			}
			data, err := os.ReadFile(tmpDir + "/test")
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, []byte(tC.data), []byte(data))
		})
	}
}
