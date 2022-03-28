package core

import (
	"encoding/json"
	"fmt"
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
	Hash string      `json:"hash,omitempty"`
	Link string      `json:"link,omitempty"`
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

type PullOptions struct {
	DryRun      bool
	NoFetch     bool
	Delete      bool
	RefOrCommit *string
	FileFilter  PathFilter
}

type PathFilter func(path string) bool

type DiffOptions struct {
	LeftRef       string
	LeftCommit    *Commit
	RightRef      string
	RightCommit   *Commit
	AddFilter     PathFilter
	ChangeFilter  PathFilter
	DeleteFilter  PathFilter
	IncludeFilter PathFilter
	NoDelete      bool
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
	Hash    string
	Link    string
	Path    string
	Size    int64
	Mode    fs.FileMode
	OldPath string
	OldLink string
	OldHash string
	OldSize int64
	OldMode fs.FileMode
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
	info, err := os.Lstat(fullPath)
	if err != nil {
		return BlobMetaData{}, err
	}

	if info.Mode()&fs.ModeSymlink != 0 {
		link, err := os.Readlink(fullPath)
		if err != nil {
			return BlobMetaData{}, err
		}

		return BlobMetaData{
			Path: path,
			Link: link,
			Mode: 0,
		}, nil
	} else if info.Mode().IsRegular() {
		hash, _ := Sha1SumFromFile(fullPath)
		return BlobMetaData{
			Path: path,
			Hash: hash,
			Mode: info.Mode().Perm(),
			Size: info.Size(),
		}, nil
	} else {
		fmt.Printf("%s %s\n", info.Mode(), info.Name())
		return BlobMetaData{}, os.ErrInvalid
	}
}

func MakeCommitMetadata(commit *Commit) ([]byte, string) {
	jsondata, _ := json.Marshal(commit)
	hash := Sha1Sum(jsondata)
	return jsondata, hash
}
