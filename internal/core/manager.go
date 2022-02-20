package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/infuseai/art/internal/executor"
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
	repoPath := MakeObjectPath(metadata.Hash)
	_, err := mngr.repo.Stat(repoPath)
	if err == nil {
		log.Printf("skip:   %s\n", metadata.Path)
		return nil
	} else {
		log.Printf("upload: %s\n", metadata.Path)
	}

	blobPath := filepath.Join(mngr.baseDir, metadata.Path)
	err = mngr.repo.Upload(blobPath, repoPath)
	return err
}

func (mngr *ArtifactMangager) DownloadBlob(metadata BlobMetaData) error {
	hash, err := Sha1SumFromFile(path.Join(mngr.baseDir, metadata.Path))
	if err == nil && hash == metadata.Hash {
		log.Printf("Skip:     %s\n", metadata.Path)
		return nil
	} else {
		log.Printf("download: %s\n", metadata.Path)
	}

	blobPath := filepath.Join(mngr.baseDir, metadata.Path)
	repoPath := MakeObjectPath(metadata.Hash)
	err = mngr.repo.Download(repoPath, blobPath)
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

func (mngr *ArtifactMangager) DeleteRef(ref string) error {
	refPath := MakeRefPath(ref)
	localPath := path.Join(mngr.metadataDir, refPath)

	err := deleteFile(localPath)
	if err != nil {
		return err
	}

	err = mngr.repo.Delete(refPath)
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

func (mngr *ArtifactMangager) FindCommitOrReference(refOrCommit string) (string, error) {
	var refPath string
	if refOrCommit == RefLatest {
		refPath = MakeRefPath(RefLatest)
	} else {
		refPath = MakeTagPath(refOrCommit)
	}

	data, err := readFile(path.Join(mngr.metadataDir, refPath))
	if err == nil {
		return string(data), nil
	}

	if len(refOrCommit) >= 4 {
		dirEntries, err := ioutil.ReadDir(path.Join(mngr.metadataDir, "commits"))
		candidates := []string{}

		if err != nil {
			return "", err
		}

		for _, entry := range dirEntries {
			if entry.IsDir() {
				continue
			}

			if strings.HasPrefix(entry.Name(), refOrCommit) {
				candidates = append(candidates, entry.Name())
			}
		}

		if len(candidates) == 1 {
			return candidates[0], nil
		}
	}

	return "", ErrReferenceNotFound
}

// Fetch downloads all the metadata from repository
func (mngr *ArtifactMangager) Fetch() error {
	// fetch latest
	mngr.GetRef(RefLatest)

	// fetch tags
	tagEntries, err := mngr.repo.List("refs/tags")
	if err != nil {
		return err
	}

	for _, entry := range tagEntries {
		mngr.GetRef("tags/" + entry.Name())
	}

	// fetch commmits
	commitEntries, err := mngr.repo.List("commits")
	if err != nil {
		return err
	}

	for _, entry := range commitEntries {
		_, err := mngr.GetCommit(entry.Name())
		if err != nil {
			return err
		}
	}

	return nil
}

func (mngr *ArtifactMangager) Push(option PushOptions) error {
	parent, err := mngr.GetRef(RefLatest)
	if err != nil {
		parent = ""
	}

	commit, err := mngr.MakeLocalCommit(parent, option.Message)
	if err != nil {
		return err
	}

	for _, metadata := range commit.Blobs {
		err := mngr.UploadBlob(metadata)
		if err != nil {
			log.Fatalf("cannot upload blob: %s\n", metadata.Path)
			break
		}
	}

	_, hash := MakeCommitMetadata(commit)
	mngr.Commit(*commit)

	mngr.AddRef(RefLatest, hash)
	if option.Tag != nil {
		mngr.AddTag(hash, *option.Tag)
	}

	return nil
}

func (mngr *ArtifactMangager) MakeLocalCommit(parent string, message *string) (*Commit, error) {
	baseDir := mngr.baseDir
	commit := Commit{
		CreatedAt: time.Now(),
		Parent:    parent,
		Message:   message,
		Blobs:     make([]BlobMetaData, 0),
	}

	tasks := []executor.TaskFunc{}
	mutex := sync.Mutex{}
	filepath.Walk(baseDir, func(absPath string, info fs.FileInfo, err error) error {
		task := func(ctx context.Context) error {
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

			mutex.Lock()
			commit.Blobs = append(commit.Blobs, metadata)
			mutex.Unlock()

			return nil
		}
		tasks = append(tasks, task)
		return nil
	})

	err := executor.ExecuteAll(0, tasks...)
	if err != nil {
		return nil, err
	}

	return &commit, nil
}

func (mngr *ArtifactMangager) Pull(options PullOptions) error {
	var err error
	if options.Fetch {
		err = mngr.Fetch()
		if err != nil {
			return err
		}
	}

	refOrCommit := RefLatest
	if options.RefOrCommit != nil {
		refOrCommit = *options.RefOrCommit
	}

	commitHash, err := mngr.FindCommitOrReference(refOrCommit)
	if err != nil {
		var refPath string
		if refOrCommit == RefLatest {
			refPath = RefLatest
		} else {
			refPath = MakeTagPath(refOrCommit)
		}
		commitHash, err = mngr.GetRef(refPath)
		if err != nil {
			return ErrReferenceNotFound
		}
	}

	commit, err := mngr.GetCommit(commitHash)
	if err != nil {
		return err
	}
	for _, blob := range commit.Blobs {
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

func (mngr *ArtifactMangager) ListTags() error {
	err := mngr.Fetch()
	if err != nil {
		return err
	}

	dirEntries, err := ioutil.ReadDir(path.Join(mngr.metadataDir, "refs/tags"))
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}

		fmt.Println(entry.Name())
	}
	return nil
}

func (mngr *ArtifactMangager) AddTag(refOrCommit, tag string) error {
	if tag == RefLatest {
		return errors.New("latest cannot be a tag")
	}

	commitHash, err := mngr.FindCommitOrReference(refOrCommit)
	if err != nil {
		return err
	}

	err = mngr.AddRef("tags/"+tag, commitHash)
	if err != nil {
		return err
	}

	return nil
}

func (mngr *ArtifactMangager) DeleteTag(tag string) error {
	if tag == RefLatest {
		return errors.New("latest cannot be a tag")
	}

	err := mngr.DeleteRef("tags/" + tag)
	if err != nil {
		return err
	}

	return nil
}

func (mngr *ArtifactMangager) List(refOrCommit string) error {
	commitHash, err := mngr.FindCommitOrReference(refOrCommit)
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

func (mngr *ArtifactMangager) Diff(leftRef, rightRef string) error {
	type DiffEntry struct {
		left  *BlobMetaData
		right *BlobMetaData
	}
	entries := map[string]DiffEntry{}

	var commitHash string
	var commit *Commit
	var err error
	// left
	if leftRef == RefLocal {
		commit, err = mngr.MakeLocalCommit("", nil)
		if err != nil {
			return err
		}
	} else {
		commitHash, err = mngr.FindCommitOrReference(leftRef)
		if err != nil {
			return err
		}

		commit, err = mngr.GetCommit(commitHash)
		if err != nil {
			return err
		}
	}
	for i, blob := range commit.Blobs {
		entry := entries[blob.Path]
		entry.left = &commit.Blobs[i]
		entries[blob.Path] = entry
	}

	// right
	commitHash, err = mngr.FindCommitOrReference(rightRef)
	if err != nil {
		return err
	}

	commit, err = mngr.GetCommit(commitHash)
	if err != nil {
		return err
	}

	for i, blob := range commit.Blobs {
		entry := entries[blob.Path]
		entry.right = &commit.Blobs[i]
		entries[blob.Path] = entry
	}

	// diff
	paths := []string{}
	for path, _ := range entries {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	for _, path := range paths {
		entry := entries[path]
		if entry.left != nil && entry.right == nil {
			color.HiGreen(fmt.Sprintf("+ %s\n", entry.left.Path))
		} else if entry.left == nil && entry.right != nil {
			color.HiRed(fmt.Sprintf("- %s\n", entry.right.Path))
		} else if entry.left.Hash != entry.right.Hash {
			color.HiYellow(fmt.Sprintf("! %s\n", entry.left.Path))
		} else {
			fmt.Printf("= %s\n", entry.left.Path)
		}
	}

	return nil
}

func (mngr *ArtifactMangager) Log(refOrCommit string) error {
	err := mngr.Fetch()
	if err != nil {
		return err
	}

	type RefEntry struct {
		refType string
		ref     string
	}

	commitIndex := map[string][]RefEntry{}

	// get latest
	latestPath := path.Join(mngr.metadataDir, "refs/latest")
	data, err := readFile(latestPath)
	if err == nil {
		commitHash := string(data)
		commitIndex[commitHash] = []RefEntry{{
			refType: RefLatest,
			ref:     RefLatest,
		}}
	} else {
		return err
	}

	// get reference
	refTagsDir := path.Join(mngr.metadataDir, "refs/tags")
	dirEntries, err := ioutil.ReadDir(refTagsDir)
	if err == nil {
		for _, entry := range dirEntries {
			ref := entry.Name()
			if entry.IsDir() {
				continue
			}

			data, err := readFile(path.Join(refTagsDir, entry.Name()))
			if err != nil {
				return err
			}

			commitHash := string(data)
			refEntry := RefEntry{
				refType: "tag",
				ref:     ref,
			}
			if commitIndex[commitHash] != nil {
				commitIndex[commitHash] = append(commitIndex[commitHash], refEntry)
			} else {
				commitIndex[commitHash] = []RefEntry{refEntry}
			}
		}
	}

	// log from refOrCommit
	commitHash, err := mngr.FindCommitOrReference(refOrCommit)
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

		color.Set(color.FgYellow)
		fmt.Printf("%s ", commitHash[:8])
		color.Set(color.FgHiBlack)
		fmt.Printf("%s ", createdAt)

		if commitIndex[commitHash] != nil {
			first := true
			color.Set(color.FgYellow)
			fmt.Print("(")
			for _, refEntry := range commitIndex[commitHash] {
				if !first {
					color.Set(color.FgYellow)
					fmt.Print(", ")
				}

				if refEntry.refType == RefLatest {
					color.Set(color.FgHiGreen)
				} else {
					color.Set(color.FgHiRed)
				}

				fmt.Print(refEntry.ref)
				first = false
			}
			color.Set(color.FgYellow)
			fmt.Print(") ")
		}

		color.Set(color.FgHiWhite)
		fmt.Printf("%s\n", message)
		color.Unset()

		commitHash = commit.Parent
		count++
	}

	return nil
}
