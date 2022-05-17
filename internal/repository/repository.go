package repository

import (
	neturl "net/url"
	"os"
	"path/filepath"
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

type RepoParseResult struct {
	Repo   string
	scheme string
	host   string
	path   string
}

func ParseRepo(repo string) (RepoParseResult, error) {
	var result RepoParseResult

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

		result.Repo = repo
		result.scheme = url.Scheme
		result.host = url.Host
		result.path = url.Path
	} else {
		i := strings.Index(repo, ":")
		if i > 0 {
			result.Repo = repo
			result.scheme = "ssh"
			result.host = repo[0:i]
			result.path = repo[i+1:]
		} else {
			cwd, err := os.Getwd()
			if err != nil {
				return result, err
			}
			if !strings.HasPrefix(repo, "/") {
				repo, err = filepath.Abs(filepath.Join(cwd, repo))
				if err != nil {
					return result, err
				}
			}

			result.Repo = repo
			result.scheme = "file"
			result.host = ""
			result.path = repo
		}
	}

	return result, nil
}

func ParseRepoName(result RepoParseResult) (string, error) {
	if result.scheme == "ssh" {
		name := filepath.Base(result.path)
		if name == "/" {
			return result.host, nil
		}
		return name, nil
	} else {
		url, err := neturl.Parse(result.Repo)
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
}

func NewRepository(result RepoParseResult) (Repository, error) {
	repo := result.Repo
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
	case "http":
		return NewHttpRepository(repo)
	case "https":
		if IsAzureStorageUrl(repo) {
			return NewAzureBlobRepository(repo)
		} else {
			return NewHttpRepository(repo)
		}
	default:
		return nil, UnsupportedRepositoryError{
			Message: "unsupported repository",
		}
	}
}
