package repository

import (
	"io/fs"
	neturl "net/url"
	"path/filepath"
	"strings"

	"github.com/infuseai/artiv/internal/meter"
)

type (
	ListEntry fs.DirEntry
	FileInfo  fs.FileInfo
)

type Repository interface {
	Upload(localPath, repoPath string, meter *meter.Meter) error
	Download(repoPath, localPath string, meter *meter.Meter) error
	Delete(repoPath string) error
	Stat(repoPath string) (FileInfo, error)
	List(repoPath string) ([]ListEntry, error)
}

func NewRepository(repo string) (Repository, error) {
	if strings.HasPrefix(repo, "/") {
		repo = "file://" + repo
	}

	url, err := neturl.Parse(repo)
	if err != nil {
		return nil, err
	}

	if url.Scheme == "" {
		return nil, UnsupportedRepositoryError{
			Message: "unsupported repository. Relative path is not allowed as a repository path",
		}
	}

	switch url.Scheme {
	case "file":
		return NewLocalFileSystemRepository(url.Path)
	case "s3":
		return NewS3Repository(url.Host, url.Path)
	case "http", "https":
		return NewHttpRepository(repo)
	default:
		return nil, UnsupportedRepositoryError{
			Message: "unsupported repository",
		}
	}
}

func ValidateRepository(repo string) error {
	if strings.HasPrefix(repo, "/") {
		repo = "file://" + repo
	}

	url, err := neturl.Parse(repo)
	if err != nil {
		return err
	}

	if url.Scheme == "" {
		return UnsupportedRepositoryError{
			Message: "unsupported repository. Relative path is not allowed as a repository path",
		}
	}

	switch url.Scheme {
	case "file", "s3", "http":
		return nil
	default:
		return UnsupportedRepositoryError{
			Message: "unsupported repository",
		}
	}
}

func ParseRepositoryName(repo string) (string, error) {
	// local repo
	if strings.HasPrefix(repo, "/") {
		return filepath.Base(repo), nil
	}

	// s3
	url, err := neturl.Parse(repo)
	if err != nil {
		return "", err
	}

	if url.Path == "" {
		return url.Hostname(), nil
	}

	name := filepath.Base(url.Path)
	if name == "/" {
		return url.Hostname(), nil
	}

	return name, nil
}
