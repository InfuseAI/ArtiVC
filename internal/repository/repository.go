package repository

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/infuseai/art/internal/core"
)

func GenObjectPath(prefix string, hash string) string {
	return fmt.Sprintf("%s/objects/%s/%s", prefix, hash[:2], hash[2:])
}

func GenCommitPath(prefix string, hash string) string {
	return fmt.Sprintf("%s/commits/%s", prefix, hash)
}

func GenRefPath(prefix string, ref string) string {
	return fmt.Sprintf("%s/refs/%s", prefix, ref)
}

type Repository interface {
	UploadBlob(metadata core.BlobMetaData) error
	DownloadBlob(metadata core.BlobMetaData) error
	Commit(commit core.Commit) error
	AddRef(ref string, commitHash string) error
}

// Local Filesystem
type LocalFileSystemRepository struct {
	BaseDir string
	RepoDir string
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	err = os.MkdirAll(filepath.Dir(dst), fs.ModePerm)
	if err != nil {
		return 0, err
	}

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func writeFile(content []byte, dst string) error {
	err := os.MkdirAll(filepath.Dir(dst), fs.ModePerm)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dst, content, 0o644)
	return err
}

func (repo *LocalFileSystemRepository) UploadBlob(metadata core.BlobMetaData) error {
	soureFile := filepath.Join(repo.BaseDir, metadata.Path)
	destFile := GenObjectPath(repo.RepoDir, metadata.Hash)
	_, error := copy(soureFile, destFile)
	return error
}

func (repo *LocalFileSystemRepository) DownloadBlob(metadata core.BlobMetaData) error {
	panic("not implemented") // TODO: Implement
}

func (repo *LocalFileSystemRepository) Commit(commit core.Commit) error {
	content, hash := core.MakeCommitMetadata(&commit)
	destFile := GenCommitPath(repo.RepoDir, hash)
	err := writeFile(content, destFile)
	return err
}

func (repo *LocalFileSystemRepository) AddRef(ref string, commitHash string) error {
	destFile := GenRefPath(repo.RepoDir, ref)
	err := writeFile([]byte(commitHash), destFile)
	return err
}
