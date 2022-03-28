package repository

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
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

	return &RcloneFileInfo{
		Name_:  filepath.Base(repoPath),
		IsDir_: false,
	}, nil
}

func (repo *RcloneRepository) List(repoPath string) ([]FileInfo, error) {
	var out bytes.Buffer
	cmd := exec.Command("rclone", "lsjson", repo.remotePath(repoPath))
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	var rcloneEntries []RcloneFileInfo
	err = json.Unmarshal(out.Bytes(), &rcloneEntries)
	if err != nil {
		return nil, err
	}

	entries := make([]FileInfo, 0)
	for _, entry := range rcloneEntries {
		entries = append(entries, &entry)
	}
	return entries, nil
}

func (repo *RcloneRepository) remotePath(repoPath string) string {
	path := filepath.Join(repo.BaseDir, repoPath)
	return repo.Remote + ":" + path
}

type RcloneFileInfo struct {
	Name_  string `json:"Name"`
	IsDir_ bool   `json:"IsDir"`
}

func (e *RcloneFileInfo) Name() string {
	return e.Name_
}

func (e *RcloneFileInfo) IsDir() bool {
	return e.IsDir_
}
