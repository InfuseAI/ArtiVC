package repository

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/infuseai/art/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestGenObjectPath(t *testing.T) {
	path := GenObjectPath("/foo/bar", "da39a3ee5e6b4b0d3255bfef95601890afd80709")
	assert.Equal(t, "/foo/bar/objects/da/39a3ee5e6b4b0d3255bfef95601890afd80709", path)
}

func TestLocalFileSystemUpload(t *testing.T) {
	baseDir := t.TempDir() + "/repo"
	repoDir := t.TempDir() + "/data"
	path := "test"
	content := "test-data"
	repo := LocalFileSystemRepository{
		BaseDir: baseDir,
		RepoDir: repoDir,
	}

	// prepare data
	fullPath := filepath.Join(repo.BaseDir, path)
	os.MkdirAll(filepath.Dir(fullPath), fs.ModePerm)
	if err := ioutil.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		assert.Fail(t, "cannot write file")
		return
	}

	// upload
	// get the size
	metaData, err := core.MakeBlobMetadata(baseDir, path)
	if err != nil {
		assert.Fail(t, "cannot create metedata")
		return
	}

	err = repo.UploadBlob(metaData)
	if err != nil {
		assert.Fail(t, "cannot create metedata")
		return
	}

	destFile := GenObjectPath(repo.RepoDir, metaData.Hash)
	data, _ := ioutil.ReadFile(destFile)
	assert.Equal(t, content, string(data))
}
