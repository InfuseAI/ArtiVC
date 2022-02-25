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
	DryRun  bool
	Message *string
	Tag     *string
}

type ChangeMode int

const (
	ChangeModeNone = iota
	ChangeModeMerge
	ChangeModeSync
)

type PullOptions struct {
	DryRun      bool
	Fetch       bool
	Diff        bool
	Mode        ChangeMode
	RefOrCommit *string
}

type DiffOptions struct {
	Mode        ChangeMode
	LeftRef     string
	LeftCommit  *Commit
	RightRef    string
	RightCommit *Commit
}

type DiffType int

const (
	DiffTypeAdd DiffType = iota
	DiffTypeDelete
	DiffTypeChange
	DiffTypeRename
)

type DiffRecord struct {
	Type    DiffType
	Path    string
	NewPath string
}

type DiffResult struct {
	Records []DiffRecord
}

type BlobDownloadResult struct {
	// File not changed. Skip the download
	Skip bool
}

type BlobUploadResult struct {
	// Blob exists in ther repo. Skip the upload
	Skip bool
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
