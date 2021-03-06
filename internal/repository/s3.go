package repository

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Repository struct {
	Bucket   string
	BasePath string
	client   *s3.Client
}

func NewS3Repository(bucket, basePath string) (*S3Repository, error) {
	basePath = strings.TrimPrefix(basePath, "/")

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	client := s3.NewFromConfig(cfg)

	return &S3Repository{
		Bucket:   bucket,
		BasePath: basePath,
		client:   client,
	}, nil
}

func (repo *S3Repository) Upload(localPath, repoPath string, m *Meter) error {
	// Reference the code to show the progress when uploading
	// https://github.com/aws/aws-sdk-go/blob/main/example/service/s3/putObjectWithProcess/putObjWithProcess.go
	sourceFileStat, err := os.Stat(localPath)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", localPath)
	}

	source, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer source.Close()

	fileInfo, err := source.Stat()
	if err != nil {
		return err
	}

	reader := &progressReader{
		fp:    source,
		size:  fileInfo.Size(),
		meter: m,
	}

	key := filepath.Join(repo.BasePath, repoPath)
	input := &s3.PutObjectInput{
		Bucket: &repo.Bucket,
		Key:    &key,
		Body:   reader,
	}

	if sourceFileStat.Size() < manager.DefaultUploadPartSize {
		_, err = repo.client.PutObject(context.TODO(), input)
	} else {
		uploader := manager.NewUploader(repo.client)
		_, err = uploader.Upload(context.TODO(), input)
	}
	return err
}

func (repo *S3Repository) Download(repoPath, localPath string, m *Meter) error {
	// Reference the code to show the progress when downloading
	// https://github.com/aws/aws-sdk-go/tree/main/example/service/s3/getObjectWithProgress
	key := filepath.Join(repo.BasePath, repoPath)
	input := &s3.GetObjectInput{
		Bucket: &repo.Bucket,
		Key:    &key,
	}

	downloader := manager.NewDownloader(repo.client)

	dest, err := os.Create(localPath)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	defer dest.Close()

	writer := &progressWriter{writer: dest, meter: m}
	_, err = downloader.Download(context.TODO(), writer, input)
	return err
}

func (repo *S3Repository) Delete(repoPath string) error {
	key := filepath.Join(repo.BasePath, repoPath)
	input := &s3.DeleteObjectInput{
		Bucket: &repo.Bucket,
		Key:    &key,
	}

	_, err := repo.client.DeleteObject(context.TODO(), input)
	return err
}

func (repo *S3Repository) Stat(repoPath string) (FileInfo, error) {
	key := filepath.Join(repo.BasePath, repoPath)
	input := &s3.HeadObjectInput{
		Bucket: &repo.Bucket,
		Key:    &key,
	}
	_, err := repo.client.HeadObject(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	return &S3FileInfo{
		name: filepath.Base(repoPath),
	}, nil
}

func (repo *S3Repository) List(repoPath string) ([]FileInfo, error) {
	fullRepoPath := filepath.Join(repo.BasePath, repoPath)
	fullRepoPath = fullRepoPath + "/"
	delimeter := "/"
	input := &s3.ListObjectsV2Input{
		Bucket:    &repo.Bucket,
		Prefix:    &fullRepoPath,
		Delimiter: &delimeter,
	}
	output, err := repo.client.ListObjectsV2(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	entries := make([]FileInfo, 0)
	for _, prefix := range output.CommonPrefixes {
		fullname := *prefix.Prefix
		name := fullname[len(fullRepoPath) : len(fullname)-1]
		entry := S3FileInfo{name: name, isDir: true}
		entries = append(entries, &entry)
	}

	for _, obj := range output.Contents {
		fullname := *obj.Key
		entry := S3FileInfo{name: fullname[len(fullRepoPath):]}
		entries = append(entries, &entry)
	}
	return entries, err
}

type S3FileInfo struct {
	name  string
	isDir bool
}

func (fi *S3FileInfo) Name() string {
	return fi.name
}

func (fi *S3FileInfo) IsDir() bool {
	return fi.isDir
}

type progressReader struct {
	fp    *os.File
	size  int64
	meter *Meter
}

func (r *progressReader) Read(p []byte) (int, error) {
	read, err := r.fp.Read(p)
	if r.meter != nil {
		r.meter.AddBytes(read)
	}
	return read, err
}

func (r *progressReader) ReadAt(p []byte, off int64) (int, error) {
	n, err := r.fp.ReadAt(p, off)
	if err != nil {
		return n, err
	}

	if r.meter != nil {
		r.meter.AddBytes(n)
	}

	return n, err
}

func (r *progressReader) Seek(offset int64, whence int) (int64, error) {
	return r.fp.Seek(offset, whence)
}

type progressWriter struct {
	writer io.WriterAt
	meter  *Meter
}

func (w *progressWriter) WriteAt(p []byte, off int64) (int, error) {
	n, err := w.writer.WriteAt(p, off)
	if err != nil {
		return n, err
	}

	if w.meter != nil {
		w.meter.AddBytes(n)
	}

	return n, err
}
