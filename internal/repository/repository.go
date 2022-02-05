package repository

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func GenObjectPath(prefix string, hash string) string {
	return fmt.Sprintf("%s/objects/%s/%s", prefix, hash[:2], hash[2:])
}

type FileMetaData struct {
	Path string
	Hash string
	Size int64
}

type Repository interface {
	UploadBlob(metaData FileMetaData)
	DownloadBlob(metaData FileMetaData)
}

// Local Filesystem
type LocalFileSystemRepository struct {
	BaseDir string
	RepoDir string
}

func copy(sourceFile string, destinationFile string) {
	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	os.MkdirAll(filepath.Dir(destinationFile), os.ModePerm)
	err = ioutil.WriteFile(destinationFile, input, 0644)
	if err != nil {
		fmt.Println("Error creating", destinationFile)
		fmt.Println(err)
		return
	}
}

func (repo *LocalFileSystemRepository) UploadBlob(metaData FileMetaData) {
	soureFile := filepath.Join(repo.BaseDir, metaData.Path)
	destFile := GenObjectPath(repo.RepoDir, metaData.Hash)

	copy(soureFile, destFile)
}

func (repo *LocalFileSystemRepository) DownloadBlob(metaData FileMetaData) {
	panic("not implemented") // TODO: Implement
}

func Size(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fi.Size())
}
