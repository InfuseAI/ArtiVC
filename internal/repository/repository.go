package repository

import "io/fs"

type ListEntry fs.DirEntry

type Repository interface {
	Upload(localPath, repoPath string) error
	Download(repoPath, localPath string) error
	Delete(repoPath string) error
	List(repoPath string) ([]ListEntry, error)
}

func NewRepository(repo string) Repository {
	return &LocalFileSystemRepository{
		RepoDir: repo,
	}
}
