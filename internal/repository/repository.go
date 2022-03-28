package repository

import (
	neturl "net/url"
	"strings"
)

type FileInfo interface {
	Name() string
	IsDir() bool
}

type Repository interface {
	Upload(localPath, repoPath string, meter *Meter) error
	Download(repoPath, localPath string, meter *Meter) error
	Delete(repoPath string) error
	Stat(repoPath string) (FileInfo, error)
	List(repoPath string) ([]FileInfo, error)
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
	case "rclone":
		return NewRcloneRepository(url.Host, url.Path)
	case "ssh":
		return NewSSHRepository(url.Host, url.Path)
	case "http", "https":
		return NewHttpRepository(repo)
	default:
		return nil, UnsupportedRepositoryError{
			Message: "unsupported repository",
		}
	}
}
