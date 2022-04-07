package repository

import (
	neturl "net/url"
	"strings"
)

type FileInfo interface {
	Name() string
	IsDir() bool
}

type SimpleFileInfo struct {
	name  string
	isDir bool
}

func (fi *SimpleFileInfo) Name() string {
	return fi.name
}

func (fi *SimpleFileInfo) IsDir() bool {
	return fi.isDir
}

type Repository interface {
	Upload(localPath, repoPath string, meter *Meter) error
	Download(repoPath, localPath string, meter *Meter) error
	Delete(repoPath string) error
	Stat(repoPath string) (FileInfo, error)
	List(repoPath string) ([]FileInfo, error)
}

type repoParseResult struct {
	scheme string
	host   string
	path   string
}

func parseRepo(repo string) (repoParseResult, error) {
	var result repoParseResult

	if strings.Contains(repo, "://") {
		url, err := neturl.Parse(repo)
		if err != nil {
			return result, err
		}

		if url.Scheme == "" {
			return result, UnsupportedRepositoryError{
				Message: "unsupported repository. Relative path is not allowed as a repository path",
			}
		}

		result.scheme = url.Scheme
		result.host = url.Host
		result.path = url.Path
	} else {
		i := strings.Index(repo, ":")
		if i > 0 {
			result.scheme = "ssh"
			result.host = repo[0:i]
			result.path = repo[i+1:]
		} else {
			result.scheme = "file"
			result.host = ""
			result.path = repo
		}
	}

	return result, nil
}

func NewRepository(repo string) (Repository, error) {
	result, err := parseRepo(repo)
	if err != nil {
		return nil, err
	}

	host := result.host
	path := result.path

	switch result.scheme {
	case "file":
		return NewLocalFileSystemRepository(path)
	case "s3":
		return NewS3Repository(host, path)
	case "gs":
		return NewGCSRepository(host, path)
	case "rclone":
		return NewRcloneRepository(host, path)
	case "ssh":
		return NewSSHRepository(host, path)
	case "http", "https":
		return NewHttpRepository(repo)
	default:
		return nil, UnsupportedRepositoryError{
			Message: "unsupported repository",
		}
	}
}
