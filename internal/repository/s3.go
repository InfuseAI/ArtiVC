package repository

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/infuseai/artiv/internal/meter"
)

type S3Repository struct {
	Bucket   string
	BasePath string
	client   *s3.Client
}

func NewS3Repository(bucket, basePath string) (*S3Repository, error) {
	basePath = strings.TrimPrefix(basePath, "/")

	return &S3Repository{
		Bucket:   bucket,
		BasePath: basePath,
	}, nil
}

func (repo *S3Repository) Upload(localPath, repoPath string, m *meter.Meter) error {
	if repo.client == nil {
		repo.init()
	}

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
		fp:      source,
		size:    fileInfo.Size(),
		signMap: map[int64]struct{}{},
		meter:   m,
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

func (repo *S3Repository) Download(repoPath, localPath string, m *meter.Meter) error {
	if repo.client == nil {
		repo.init()
	}

	key := filepath.Join(repo.BasePath, repoPath)
	input := &s3.GetObjectInput{
		Bucket: &repo.Bucket,
		Key:    &key,
	}

	_, err := repo.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: &repo.Bucket,
		Key:    &key,
	})
	if err != nil {
		return err
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
	if repo.client == nil {
		repo.init()
	}

	key := filepath.Join(repo.BasePath, repoPath)
	input := &s3.DeleteObjectInput{
		Bucket: &repo.Bucket,
		Key:    &key,
	}

	_, err := repo.client.DeleteObject(context.TODO(), input)
	return err
}

func (repo *S3Repository) Stat(repoPath string) (FileInfo, error) {
	if repo.client == nil {
		repo.init()
	}

	key := filepath.Join(repo.BasePath, repoPath)
	input := &s3.HeadObjectInput{
		Bucket: &repo.Bucket,
		Key:    &key,
	}
	_, err := repo.client.HeadObject(context.TODO(), input)
	return nil, err
}

func (repo *S3Repository) List(repoPath string) ([]ListEntry, error) {
	if repo.client == nil {
		repo.init()
	}

	fullRepoPath := filepath.Join(repo.BasePath, repoPath)
	input := &s3.ListObjectsV2Input{
		Bucket: &repo.Bucket,
		Prefix: &fullRepoPath,
	}
	output, err := repo.client.ListObjectsV2(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	entries := make([]ListEntry, 0)
	for _, obj := range output.Contents {
		key := *obj.Key
		entry := S3DirEntry{name: key[len(fullRepoPath)+1:]}
		entries = append(entries, &entry)
	}
	return entries, err
}

func (repo *S3Repository) init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	repo.client = s3.NewFromConfig(cfg)
}

type S3DirEntry struct {
	name string
}

func (e *S3DirEntry) Name() string {
	return e.name
}

func (e *S3DirEntry) IsDir() bool {
	return false
}

func (e *S3DirEntry) Type() fs.FileMode {
	return os.ModePerm
}

func (e *S3DirEntry) Info() (fs.FileInfo, error) {
	return nil, nil
}

type progressReader struct {
	fp      *os.File
	size    int64
	read    int64
	signMap map[int64]struct{}
	mux     sync.Mutex
	meter   *meter.Meter
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

	r.mux.Lock()
	// Ignore the first signature call
	if _, ok := r.signMap[off]; ok {
		// Got the length have read( or means has uploaded), and you can construct your message
		r.read += int64(n)
		if r.meter != nil {
			r.meter.AddBytes(n)
		}
	} else {
		r.signMap[off] = struct{}{}
	}
	r.mux.Unlock()
	return n, err
}

func (r *progressReader) Seek(offset int64, whence int) (int64, error) {
	return r.fp.Seek(offset, whence)
}

type progressWriter struct {
	writer io.WriterAt
	meter  *meter.Meter
}

func (w *progressWriter) WriteAt(p []byte, off int64) (int, error) {
	if w.meter != nil {
		w.meter.AddBytes(len(p))
	}
	return w.writer.WriteAt(p, off)
}
