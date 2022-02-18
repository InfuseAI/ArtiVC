package repository

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Repository struct {
	Bucket string
	client *s3.Client
}

func NewS3Repository(bucket string) (*S3Repository, error) {
	return &S3Repository{
		Bucket: bucket,
	}, nil
}

func (repo *S3Repository) Upload(localPath, repoPath string) error {
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

	input := &s3.PutObjectInput{
		Bucket: &repo.Bucket,
		Key:    &repoPath,
		Body:   source,
	}

	_, err = repo.client.PutObject(context.TODO(), input)
	return err
}

func (repo *S3Repository) Download(repoPath, localPath string) error {
	if repo.client == nil {
		repo.init()
	}

	input := &s3.GetObjectInput{
		Bucket: &repo.Bucket,
		Key:    &repoPath,
	}

	output, err := repo.client.GetObject(context.TODO(), input)
	if err != nil {
		return err
	}

	dest, err := os.Create(localPath)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	defer dest.Close()
	_, err = io.Copy(dest, output.Body)

	return err
}

func (repo *S3Repository) Delete(repoPath string) error {
	if repo.client == nil {
		repo.init()
	}

	input := &s3.DeleteObjectInput{
		Bucket: &repo.Bucket,
		Key:    &repoPath,
	}

	_, err := repo.client.DeleteObject(context.TODO(), input)
	return err
}

func (repo *S3Repository) Stat(repoPath string) (FileInfo, error) {
	if repo.client == nil {
		repo.init()
	}

	input := &s3.HeadObjectInput{
		Bucket: &repo.Bucket,
		Key:    &repoPath,
	}
	_, err := repo.client.HeadObject(context.TODO(), input)
	return nil, err
}

func (repo *S3Repository) List(repoPath string) ([]ListEntry, error) {
	if repo.client == nil {
		repo.init()
	}

	input := &s3.ListObjectsV2Input{
		Bucket: &repo.Bucket,
		Prefix: &repoPath,
	}
	output, err := repo.client.ListObjectsV2(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	entries := make([]ListEntry, 0)
	for _, obj := range output.Contents {
		key := *obj.Key
		entry := S3DirEntry{name: key[len(repoPath)+1:]}
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
