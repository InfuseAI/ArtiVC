package repository

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Local Filesystem
type RcloneRepository struct {
	Remote  string
	BaseDir string
}

func NewRcloneRepository(remote, basePath string) (*RcloneRepository, error) {
	cmd := exec.Command("rclone", "version")
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	return &RcloneRepository{
		Remote:  remote,
		BaseDir: basePath,
	}, nil
}

func (repo *RcloneRepository) Upload(localPath, repoPath string, m *Meter) error {
	cmd := exec.Command("rclone", "copyto", "--no-check-dest", localPath, repo.remotePath(repoPath))
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (repo *RcloneRepository) Download(repoPath, localPath string, m *Meter) error {
	cmd := exec.Command("rclone", "copyto", "--no-check-dest", repo.remotePath(repoPath), localPath)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (repo *RcloneRepository) Delete(repoPath string) error {
	cmd := exec.Command("rclone", "deletefile", repo.remotePath(repoPath))
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (repo *RcloneRepository) Stat(repoPath string) (FileInfo, error) {
	var out bytes.Buffer
	cmd := exec.Command("rclone", "size", "--json", repo.remotePath(repoPath))
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	type RcloneSize struct {
		Count int `json:"count"`
	}

	var size RcloneSize
	json.Unmarshal(out.Bytes(), &size)

	if size.Count == 0 {
		return nil, os.ErrNotExist
	}

	return nil, nil
}

func (repo *RcloneRepository) List(repoPath string) ([]ListEntry, error) {
	var out bytes.Buffer
	cmd := exec.Command("rclone", "lsjson", repo.remotePath(repoPath))
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	var rcloneEntries []RcloneListEntry
	err = json.Unmarshal(out.Bytes(), &rcloneEntries)
	if err != nil {
		return nil, err
	}

	entries := make([]ListEntry, 0)
	for _, entry := range rcloneEntries {
		entries = append(entries, &entry)
	}
	return entries, nil
}

func (repo *RcloneRepository) remotePath(repoPath string) string {
	path := filepath.Join(repo.BaseDir, repoPath)
	return repo.Remote + ":" + path
}

type RcloneListEntry struct {
	Path    string    `json:"Path"`
	Name_   string    `json:"Name"`
	Size    uint64    `json:"Size"`
	ModTime time.Time `json:"ModTime"`
	IsDir_  bool      `json:"IsDir"`
}

func (e *RcloneListEntry) Name() string {
	return e.Name_
}

func (e *RcloneListEntry) IsDir() bool {
	return e.IsDir_
}

func (e *RcloneListEntry) Type() fs.FileMode {
	return os.ModePerm
}

func (e *RcloneListEntry) Info() (fs.FileInfo, error) {
	return nil, nil
}
