package repository

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/infuseai/artiv/internal/meter"
)

// Local Filesystem
type LocalFileSystemRepository struct {
	RepoDir string
}

func NewLocalFileSystemRepository(repoDir string) (*LocalFileSystemRepository, error) {
	stat, err := os.Stat(repoDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(repoDir, fs.ModePerm)
			if err != nil {
				return nil, errors.New("cannot make directory: " + repoDir)
			}
		} else {
			return nil, err
		}
	} else {
		if !stat.IsDir() {
			return nil, errors.New(repoDir + " is not a directory")
		}
	}

	return &LocalFileSystemRepository{
		RepoDir: repoDir,
	}, nil
}

func (repo *LocalFileSystemRepository) Upload(localPath, repoPath string, m *meter.Meter) error {
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

	// Copy from source to tmp
	tmpDir := path.Join(repo.RepoDir, "tmp")
	err = os.MkdirAll(tmpDir, fs.ModePerm)
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp(tmpDir, "*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	_, err = meter.CopyWithMeter(tmp, source, m)
	if err != nil {
		return err
	}
	err = tmp.Close()
	if err != nil {
		return err
	}

	// Move from tmp to dest
	destPath := path.Join(repo.RepoDir, repoPath)
	err = os.MkdirAll(filepath.Dir(destPath), fs.ModePerm)
	if err != nil {
		return err
	}
	err = os.Remove(destPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	err = os.Rename(tmpPath, destPath)
	if err != nil {
		return err
	}

	return nil
}

func (repo *LocalFileSystemRepository) Download(repoPath, localPath string, m *meter.Meter) error {
	srcPath := path.Join(repo.RepoDir, repoPath)
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dest, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer dest.Close()
	written, err := meter.CopyWithMeter(dest, src, m)
	if err != nil {
		return err
	}

	if written == 0 {
		err = os.Truncate(localPath, 0)
	}

	return err
}

func (repo *LocalFileSystemRepository) Delete(repoPath string) error {
	filePath := path.Join(repo.RepoDir, repoPath)
	return os.Remove(filePath)
}

func (repo *LocalFileSystemRepository) Stat(repoPath string) (FileInfo, error) {
	filePath := path.Join(repo.RepoDir, repoPath)
	return os.Stat(filePath)
}

func (repo *LocalFileSystemRepository) List(repoPath string) ([]ListEntry, error) {
	dir := path.Join(repo.RepoDir, repoPath)
	fs, err := os.ReadDir(dir)
	if err != nil {
		return []ListEntry{}, nil
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
