package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/infuseai/art/internal/repository"
)

type ArtifactMangager struct {

	// local
	baseDir string // the workspace base dir

	/* the path to store the metadata of the repository.
	If it is in an art workspace, the path should be "${baseDir}/.art"*/
	metadataDir string

	// repository
	repo repository.Repository
}

func NewArtifactManager(config ArtConfig) (*ArtifactMangager, error) {
	// init the workspace path
	baseDir := config.BaseDir
	if baseDir == "" {
		return nil, errors.New("no repository specified")
	}

	// init the metadata path
	metadataDir := config.MetadataDir
	if baseDir == "" {
		return nil, errors.New("no metadata specified")
	}

	// init the repository
	repoStr := config.RepoUrl()
	if repoStr == "" {
		return nil, errors.New("no repository specified")
	}
	repo := repository.NewRepository(repoStr)

	return &ArtifactMangager{baseDir: baseDir, repo: repo, metadataDir: metadataDir}, nil
}

func (mngr *ArtifactMangager) UploadBlob(metadata BlobMetaData) error {
	blobPath := filepath.Join(mngr.baseDir, metadata.Path)
	repoPath := MakeObjectPath(metadata.Hash)
	err := mngr.repo.Upload(blobPath, repoPath)
	return err
}

func (mngr *ArtifactMangager) DownloadBlob(metadata BlobMetaData) error {
	blobPath := filepath.Join(mngr.baseDir, metadata.Path)
	repoPath := MakeObjectPath(metadata.Hash)
	err := mngr.repo.Download(repoPath, blobPath)
	return err
}

func (mngr *ArtifactMangager) Commit(commit Commit) error {
	content, hash := MakeCommitMetadata(&commit)
	commitPath := MakeCommitPath(hash)
	localPath := path.Join(mngr.metadataDir, commitPath)
	err := writeFile(content, localPath)
	if err != nil {
		return err
	}

	err = mngr.repo.Upload(localPath, commitPath)
	if err != nil {
		return err
	}

	return err
}

func (mngr *ArtifactMangager) AddRef(ref string, commit string) error {
	refPath := MakeRefPath(ref)
	localPath := path.Join(mngr.metadataDir, refPath)
	err := writeFile([]byte(commit), localPath)
	if err != nil {
		return err
	}

	err = mngr.repo.Upload(localPath, refPath)
	if err != nil {
		return err
	}

	return err
}

func (mngr *ArtifactMangager) GetRef(ref string) (string, error) {
	refPath := MakeRefPath(ref)
	localPath := path.Join(mngr.metadataDir, refPath)

	err := mkdirsForFile(localPath)
	if err != nil {
		return "", err
	}

	err = mngr.repo.Download(refPath, localPath)
	if err != nil {
		return "", err
	}

	data, err := readFile(path.Join(mngr.metadataDir, refPath))
	if err != nil {
		return "", err
	}

	hash := string(data)
	if hash == "" {
		return "", fmt.Errorf("the reference %s is empty", ref)
	}

	return hash, nil
}

func (mngr *ArtifactMangager) GetCommit(hash string) (*Commit, error) {
	commitPath := MakeCommitPath(hash)
	localPath := path.Join(mngr.metadataDir, commitPath)

	err := mkdirsForFile(localPath)
	if err != nil {
		return nil, err
	}

	err = mngr.repo.Download(commitPath, localPath)
	if err != nil {
		return nil, err
	}

	data, err := readFile(path.Join(mngr.metadataDir, commitPath))
	if err != nil {
		return nil, err
	}

	var commit Commit
	err = json.Unmarshal(data, &commit)
	if err != nil {
		return nil, err
	}

	return &commit, nil
}

func (mngr *ArtifactMangager) Push(option PushOption) error {
	ref := "latest"
	commitHash, err := mngr.GetRef(ref)
	if err != nil {
		return err
	}

	baseDir := mngr.baseDir
	commit := Commit{
		CreatedAt: time.Now(),
		Parent:    commitHash,
		Message:   option.Message,
		Blobs:     make([]BlobMetaData, 0),
	}

	filepath.Walk(baseDir, func(absPath string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", absPath, err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		path := absPath[len(baseDir)+1:]
		if strings.HasPrefix(path, ".art") {
			return nil
		}

		metadata, err := MakeBlobMetadata(baseDir, path)
		if err != nil {
			log.Fatalf("cannot make metadata: %s", path)
			return err
		}

		commit.Blobs = append(commit.Blobs, metadata)
		return nil
	})

	for _, metadata := range commit.Blobs {
		log.Printf("upload %s\n", metadata.Path)
		err := mngr.UploadBlob(metadata)
		if err != nil {
			log.Fatalf("cannot upload blob: %s\n", metadata.Path)
			break
		}
	}

	_, hash := MakeCommitMetadata(&commit)
	mngr.Commit(commit)
	mngr.AddRef("latest", hash)

	return nil
}

func (mngr *ArtifactMangager) Pull() error {
	ref := "latest"
	commitHash, err := mngr.GetRef(ref)
	if err != nil {
		return err
	}

	commit, err := mngr.GetCommit(commitHash)
	if err != nil {
		return err
	}
	for _, blob := range commit.Blobs {
		log.Printf("download %s\n", blob.Path)

		err := mkdirsForFile(path.Join(mngr.baseDir, blob.Path))
		if err != nil {
			return err
		}

		err = mngr.DownloadBlob(blob)
		if err != nil {
			return err
		}
	}

	return nil
}

func (mngr *ArtifactMangager) List(ref string) error {
	commitHash, err := mngr.GetRef(ref)
	if err != nil {
		return err
	}

	commit, err := mngr.GetCommit(commitHash)
	if err != nil {
		return err
	}
	for _, blob := range commit.Blobs {
		fmt.Println(blob.Path)
	}

	return nil
}

func (mngr *ArtifactMangager) Log(ref string) error {
	commitHash, err := mngr.GetRef(ref)
	if err != nil {
		return err
	}

	for count := 0; commitHash != "" && count < 1000; count++ {
		commit, err := mngr.GetCommit(commitHash)
		if err != nil {
			return err
		}

		message := ""
		if commit.Message != nil {
			message = *commit.Message
		}

		createdAt := commit.CreatedAt.Format("2006-01-02 15:04 -0700")

		// fmt.Printf("%s %v %s\n", commitHash[:8], createdAt, message)
		color.Set(color.FgYellow)
		fmt.Printf("%s ", commitHash[:8])
		color.Set(color.FgHiBlack)
		fmt.Printf("%s ", createdAt)
		color.Set(color.FgHiWhite)
		fmt.Printf("%s\n", message)
		color.Unset()

		commitHash = commit.Parent
		count++
	}

	return nil
}
