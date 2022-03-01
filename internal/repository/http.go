package repository

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/infuseai/artiv/internal/meter"
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

func (repo *HttpRepository) Upload(localPath, repoPath string, meter *meter.Meter) error {
	return errors.New("Upload is not supported in Http repository")
}

func (repo *HttpRepository) Download(repoPath, localPath string, m *meter.Meter) error {
	filePath, err := getFilePath(repo.RepoUrl, repoPath)
	if err != nil {
		return err
	}

	res, err := http.Get(filePath)
	if err != nil {
		return err
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

	_, err = meter.CopyWithMeter(outputFile, res.Body, m)
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

	fileSize, err := strconv.ParseInt(res.Header["Content-Length"][0], 10, 64)
	if err != nil {
		return nil, err
	}

	modifiedTime, err := time.Parse(time.RFC1123, res.Header["Last-Modified"][0])
	if err != nil {
		return nil, err
	}

	info := &HttpFileInfo{
		name:         repoPath,
		size:         fileSize,
		modifiedTime: modifiedTime,
	}

	return info, nil
}

func (repo *HttpRepository) List(repoPath string) ([]ListEntry, error) {
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
	name         string
	size         int64
	modifiedTime time.Time
}

func (info *HttpFileInfo) Name() string {
	return info.name
}

func (info *HttpFileInfo) Size() int64 {
	return info.size
}

func (info *HttpFileInfo) Mode() fs.FileMode {
	return os.ModePerm
}

func (info *HttpFileInfo) ModTime() time.Time {
	return info.modifiedTime
}

func (info *HttpFileInfo) IsDir() bool {
	return false
}

func (info *HttpFileInfo) Sys() interface{} {
	return nil
}
