package repository

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

// Local Filesystem
type GCSRepository struct {
	Bucket   string
	BasePath string
	Client   *storage.Client
}

func NewGCSRepository(bucket, basePath string) (*GCSRepository, error) {
	ctx := context.Background()
	basePath = strings.TrimPrefix(basePath, "/")
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return &GCSRepository{
		Bucket:   bucket,
		BasePath: basePath,
		Client:   client,
	}, nil
}

func (repo *GCSRepository) Upload(localPath, repoPath string, m *Meter) error {
	ctx := context.Background()

	// client, bucket, obj
	client := repo.Client
	bkt := client.Bucket(repo.Bucket)
	obj := bkt.Object(filepath.Join(repo.BasePath, repoPath))

	// src
	src, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer src.Close()

	// dest

	dest := obj.NewWriter(ctx)
	defer dest.Close()

	// copy
	_, err = CopyWithMeter(dest, src, m)
	if err != nil {
		return err
	}

	return nil
}

func (repo *GCSRepository) Download(repoPath, localPath string, m *Meter) error {
	ctx := context.Background()

	// client, bucket, obj
	client := repo.Client
	bkt := client.Bucket(repo.Bucket)
	obj := bkt.Object(filepath.Join(repo.BasePath, repoPath))

	// src
	src, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}
	defer src.Close()

	// dest
	dest, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	// copy
	_, err = CopyWithMeter(dest, src, m)
	if err != nil {
		return err
	}

	return nil
}

func (repo *GCSRepository) Delete(repoPath string) error {
	ctx := context.Background()

	// client, bucket, obj
	client := repo.Client
	bkt := client.Bucket(repo.Bucket)
	obj := bkt.Object(filepath.Join(repo.BasePath, repoPath))

	// delete
	err := obj.Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repo *GCSRepository) Stat(repoPath string) (FileInfo, error) {
	ctx := context.Background()

	// client, bucket, obj
	client := repo.Client
	bkt := client.Bucket(repo.Bucket)
	obj := bkt.Object(filepath.Join(repo.BasePath, repoPath))

	// get object stat
	_, err := obj.Attrs(ctx)
	if err != nil {
		return nil, err
	}

	return &GCSFileInfo{
		name:  filepath.Base(repoPath),
		isDir: false,
	}, nil
}

func (repo *GCSRepository) List(repoPath string) ([]FileInfo, error) {
	ctx := context.Background()
	records := []FileInfo{}

	// client, bucket, obj
	client := repo.Client
	bkt := client.Bucket(repo.Bucket)
	prefix := filepath.Join(repo.BasePath, repoPath) + "/"
	query := &storage.Query{Prefix: prefix, Delimiter: "/"}

	it := bkt.Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return records, err
		}

		fileinfo := GCSFileInfo{}

		if attrs.Name != "" {
			fileinfo.name = attrs.Name[len(prefix):]
			fileinfo.isDir = false
		} else {
			fileinfo.name = attrs.Prefix[len(prefix) : len(attrs.Prefix)-1]
			fileinfo.isDir = true
		}
		records = append(records, &fileinfo)
	}

	return records, nil
}

type GCSFileInfo struct {
	name  string
	isDir bool
}

func (fi *GCSFileInfo) Name() string {
	return fi.name
}

func (fi *GCSFileInfo) IsDir() bool {
	return fi.isDir
}
