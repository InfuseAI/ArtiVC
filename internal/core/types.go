package core

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

const (
	RefLocal  = "__local__"
	RefLatest = "latest"
)

type BlobMetaData struct {
	Path string      `json:"path"`
	Hash string      `json:"hash"`
	Mode fs.FileMode `json:"mode"`
	Size int64       `json:"size"`
}

type Commit struct {
	CreatedAt time.Time      `json:"createdAt"`
	Parent    string         `json:"parent,omitempty"`
	Message   *string        `json:"messaage,omitempty"`
	Blobs     []BlobMetaData `json:"blobs"`
}

type PushOptions struct {
	Message *string
	Tag     *string
}

type PullOptions struct {
	Fetch bool
}

func MakeBlobMetadata(baseDir string, path string) (BlobMetaData, error) {
	fullPath := filepath.Join(baseDir, path)
	hash, _ := Sha1SumFromFile(fullPath)
	fi, err := os.Stat(fullPath)
	if err != nil {
		return BlobMetaData{}, err
	}

	metaData := BlobMetaData{
		Path: path,
		Hash: hash,
		Mode: fi.Mode(),
		Size: fi.Size(),
	}
	return metaData, nil
}

func MakeCommitMetadata(commit *Commit) ([]byte, string) {
	jsondata, _ := json.Marshal(commit)
	hash := Sha1Sum(jsondata)
	return jsondata, hash
}
