package repository

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type HttpRepository struct {
	RepoUrl string
}

func NewHttpRepository(repo string) (*HttpRepository, error) {
	if !strings.HasSuffix(repo, "/") {
		repo += "/"
	}

	return &HttpRepository{
		RepoUrl: repo,
	}, nil
}

func (repo *HttpRepository) Upload(localPath, repoPath string, meter *Meter) error {
	return errors.New("Upload is not supported in Http repository")
}

func (repo *HttpRepository) Download(repoPath, localPath string, m *Meter) error {
	filePath, err := getFilePath(repo.RepoUrl, repoPath)
	if err != nil {
		return err
	}

	res, err := http.Get(filePath)
	if err != nil {
		retry := 0
		msg := err.Error()

		for err != nil && strings.HasSuffix(msg, "connection reset by peer") && retry < 10 {
			retry++
			time.Sleep(time.Millisecond * 50 * time.Duration(retry))
			res, err = http.Get(filePath)
		}

		if err != nil {
			return err
		}
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("status code: %d", res.StatusCode)
	}

	outputFile, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = CopyWithMeter(outputFile, res.Body, m)
	return err
}

func (repo *HttpRepository) Delete(repoPath string) error {
	return errors.New("Delete is not supported in Http repository")
}

func (repo *HttpRepository) Stat(repoPath string) (FileInfo, error) {
	filePath, err := getFilePath(repo.RepoUrl, repoPath)
	if err != nil {
		return nil, err
	}

	res, err := http.Head(filePath)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code: %d", res.StatusCode)
	}

	info := &HttpFileInfo{
		name: repoPath,
	}

	return info, nil
}

func (repo *HttpRepository) List(repoPath string) ([]FileInfo, error) {
	return nil, errors.New("List is not supported in Http repository")
}

func getFilePath(repoPath, filePath string) (string, error) {
	base, err := url.Parse(repoPath)
	if err != nil {
		return "", err
	}
	path, err := url.Parse(filePath)
	if err != nil {
		return "", err
	}
	return base.ResolveReference(path).String(), nil
}

type HttpFileInfo struct {
	name string
}

func (info *HttpFileInfo) Name() string {
	return info.name
}

func (info *HttpFileInfo) IsDir() bool {
	return false
}
