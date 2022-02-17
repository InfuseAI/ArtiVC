package repository

import "io/fs"

type ListEntry fs.DirEntry
type FileInfo fs.FileInfo

type Repository interface {
	Upload(localPath, repoPath string) error
	Download(repoPath, localPath string) error
	Delete(repoPath string) error
	Stat(repoPath string) (FileInfo, error)
	List(repoPath string) ([]ListEntry, error)
}

func NewRepository(repo string) Repository {
	return &LocalFileSystemRepository{
		RepoDir: repo,
	}
}
