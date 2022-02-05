package repository

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/infuseai/art/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestGenObjectPath(t *testing.T) {
	path := GenObjectPath("/foo/bar", "da39a3ee5e6b4b0d3255bfef95601890afd80709")
	assert.Equal(t, "/foo/bar/objects/da/39a3ee5e6b4b0d3255bfef95601890afd80709", path)
}

func TestLocalFileSystemUpload(t *testing.T) {
	repo := LocalFileSystemRepository{
		BaseDir: "/Users/popcorny/art/myart",
		RepoDir: "/Users/popcorny/art/myrepo",
	}

	path := "abc"
	fullPath := filepath.Join(repo.BaseDir, path)
	hash := util.Sha1SumFromFile(fullPath)
	fi, err := os.Stat(fullPath)
	if err != nil {
		panic(err)
	}
	// get the size
	size := fi.Size()

	repo.UploadBlob(FileMetaData{
		Path: path,
		Hash: hash,
		Size: size,
	})
}
