package repository

import (
	"io/fs"
	neturl "net/url"
	"strings"
)

type ListEntry fs.DirEntry
type FileInfo fs.FileInfo

type Repository interface {
	Upload(localPath, repoPath string) error
	Download(repoPath, localPath string) error
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
	default:
		return nil, UnsupportedRepositoryError{
			Message: "unsupported repository",
		}
	}
}
