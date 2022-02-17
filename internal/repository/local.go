package repository

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

// Local Filesystem
type LocalFileSystemRepository struct {
	RepoDir string
}

func (repo *LocalFileSystemRepository) Upload(localPath, repoPath string) error {
	sourceFileStat, err := os.Stat(localPath)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", localPath)
	}

	source, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer source.Close()

	destPath := path.Join(repo.RepoDir, repoPath)
	err = os.MkdirAll(filepath.Dir(destPath), fs.ModePerm)
	if err != nil {
		return err
	}

	destination, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)

	return err
}

func (repo *LocalFileSystemRepository) Download(repoPath, localPath string) error {
	srcPath := path.Join(repo.RepoDir, repoPath)
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dest, err := os.Create(localPath)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	defer dest.Close()
	_, err = io.Copy(dest, src)

	return err
}

func (repo *LocalFileSystemRepository) Delete(repoPath string) error {
	srcPath := path.Join(repo.RepoDir, repoPath)
	return os.Remove(srcPath)
}

func (repo *LocalFileSystemRepository) List(repoPath string) ([]ListEntry, error) {
	dir := path.Join(repo.RepoDir, repoPath)
	fs, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	fs2 := []ListEntry{}

	for _, info := range fs {
		info2, ok := info.(ListEntry)
		if ok {
			fs2 = append(fs2, info2)
		}
	}
	return fs2, nil
}
