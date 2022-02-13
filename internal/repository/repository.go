package repository

type Repository interface {
	Upload(localPath, repoPath string) error
	Download(repoPath, localPath string) error
}

func NewRepository(repo string) Repository {
	return &LocalFileSystemRepository{
		RepoDir: repo,
	}
}
