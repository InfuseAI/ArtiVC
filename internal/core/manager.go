package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/infuseai/artivc/internal/executor"
	"github.com/infuseai/artivc/internal/log"
	"github.com/infuseai/artivc/internal/repository"
)

type ArtifactManager struct {

	// local
	baseDir string // the workspace base dir

	/* the path to store the metadata of the repository.
	If it is in a workspace, the path should be "${baseDir}/.avc"*/
	metadataDir string

	// repository
	repo repository.Repository
}

func NewArtifactManager(config ArtConfig) (*ArtifactManager, error) {
	// init the workspace path
	baseDir := config.BaseDir
	if baseDir == "" {
		return nil, errors.New("no workspace specified")
	}

	finfo, err := os.Stat(baseDir)
	if err == nil {
		if !finfo.IsDir() {
			return nil, errors.New(baseDir + " is not a directory")
		}
	}

	// init the metadata path
	metadataDir := config.MetadataDir
	if metadataDir == "" {
		return nil, errors.New("no metadata specified")
	}

	// init the repository
	repoStr := config.RepoUrl()
	if repoStr == "" {
		return nil, errors.New("no repository specified")
	}
	repo, err := repository.NewRepository(repoStr)
	if err != nil {
		return nil, err
	}

	return &ArtifactManager{baseDir: baseDir, repo: repo, metadataDir: metadataDir}, nil
}

func (mngr *ArtifactManager) UploadBlob(localPath, hash string, meter *repository.Meter, checkSkip bool) (BlobUploadResult, error) {
	repoPath := MakeObjectPath(hash)

	if checkSkip {
		_, err := mngr.repo.Stat(repoPath)
		if err == nil {
			log.Debugf("skip: %s\n", repoPath)
			return BlobUploadResult{Skip: true}, nil
		}
	}

	blobPath := filepath.Join(mngr.baseDir, localPath)
	err := mngr.Upload(blobPath, repoPath, meter)
	return BlobUploadResult{Skip: false}, err
}

func (mngr *ArtifactManager) Upload(localPath, repoPath string, meter *repository.Meter) error {
	log.Debugf("upload: %s -> %s\n", localPath, repoPath)

	return mngr.repo.Upload(localPath, repoPath, meter)
}

func (mngr *ArtifactManager) Download(repoPath, localPath, tmpDir string, meter *repository.Meter) error {
	log.Debugf("download: %s <- %s\n", localPath, repoPath)

	// Copy from repo to tmp
	err := os.MkdirAll(tmpDir, fs.ModePerm)
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp(tmpDir, "*")
	if err != nil {
		return err
	}
	tmp.Close()

	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	err = mngr.repo.Download(repoPath, tmpPath, meter)
	if err != nil {
		return err
	}

	// Move from tmp to local
	err = os.MkdirAll(filepath.Dir(localPath), fs.ModePerm)
	if err != nil {
		return err
	}
	err = os.Remove(localPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	err = os.Rename(tmpPath, localPath)
	if err != nil {
		return err
	}

	return nil
}

func (mngr *ArtifactManager) DownloadBlob(localPath, hash string, meter *repository.Meter) (BlobDownloadResult, error) {
	blobPath := filepath.Join(mngr.baseDir, localPath)

	err := mkdirsForFile(blobPath)
	if err != nil {
		return BlobDownloadResult{}, err
	}

	repoPath := MakeObjectPath(hash)
	tmpDir := path.Join(mngr.baseDir, ".avc", "tmp")

	err = mngr.Download(repoPath, blobPath, tmpDir, meter)
	if err != nil {
		return BlobDownloadResult{}, err
	}
	return BlobDownloadResult{Skip: false}, nil
}

func (mngr *ArtifactManager) Commit(commit Commit) error {
	content, hash := MakeCommitMetadata(&commit)
	commitPath := MakeCommitPath(hash)
	localPath := path.Join(mngr.metadataDir, commitPath)
	err := writeGzipFile(content, localPath)
	if err != nil {
		return err
	}

	err = mngr.Upload(localPath, commitPath, nil)
	if err != nil {
		return err
	}

	return err
}

func (mngr *ArtifactManager) AddRef(ref string, commit string) error {
	refPath := MakeRefPath(ref)
	localPath := path.Join(mngr.metadataDir, refPath)
	err := writeFile([]byte(commit), localPath)
	if err != nil {
		return err
	}

	err = mngr.Upload(localPath, refPath, nil)
	if err != nil {
		return err
	}

	return err
}

func (mngr *ArtifactManager) DeleteRef(ref string) error {
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

func (mngr *ArtifactManager) GetRef(ref string) (string, error) {
	refPath := MakeRefPath(ref)
	localPath := path.Join(mngr.metadataDir, refPath)

	err := mkdirsForFile(localPath)
	if err != nil {
		return "", err
	}

	tmpDir := path.Join(mngr.metadataDir, "tmp")
	err = mngr.Download(refPath, localPath, tmpDir, nil)
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

func (mngr *ArtifactManager) GetCommit(hash string) (*Commit, error) {
	commitPath := MakeCommitPath(hash)
	localPath := path.Join(mngr.metadataDir, commitPath)

	_, err := os.Stat(localPath)
	if err != nil {
		err := mkdirsForFile(localPath)
		if err != nil {
			return nil, err
		}

		tmpDir := path.Join(mngr.metadataDir, "tmp")
		err = mngr.Download(commitPath, localPath, tmpDir, nil)
		if err != nil {
			return nil, err
		}
	}

	data, err := readGzipFile(path.Join(mngr.metadataDir, commitPath))
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

func (mngr *ArtifactManager) FindCommitOrReference(refOrCommit string) (string, error) {
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
			dirEntries = []fs.FileInfo{}
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

	if refOrCommit == RefLatest {
		return "", ErrEmptyRepository
	} else {
		return "", ReferenceNotFoundError{
			Ref: refOrCommit,
		}
	}
}

// Fetch downloads all the metadata from repository
func (mngr *ArtifactManager) Fetch() error {
	log.Debugln("fetch the repository metadata")
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

func (mngr *ArtifactManager) Push(options PushOptions) error {
	parent, err := mngr.GetRef(RefLatest)
	if err != nil {
		parent = ""
	}
	checkSkip := true
	avcIgnore, err := NewAvcIgnore(mngr.baseDir)
	avcIgnoreFilter := func(path string) bool {
		return true
	}

	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		avcIgnoreFilter = func(path string) bool {
			return !avcIgnore.MatchesPath(path)
		}
	}

	commit, err := mngr.MakeWorkspaceCommit(parent, options.Message, avcIgnoreFilter)
	if err != nil {
		return err
	}

	result, err := mngr.Diff(DiffOptions{
		LeftRef:      RefLatest,
		RightCommit:  commit,
		AddFilter:    avcIgnoreFilter,
		ChangeFilter: avcIgnoreFilter,
		DeleteFilter: nil,
	})
	if err != nil {
		if err != ErrEmptyRepository {
			return err
		} else {
			checkSkip = false
			result, err = mngr.Diff(DiffOptions{
				LeftCommit:   mngr.MakeEmptyCommit(),
				RightCommit:  commit,
				AddFilter:    avcIgnoreFilter,
				ChangeFilter: avcIgnoreFilter,
				DeleteFilter: nil,
			})
			if err != nil {
				return err
			}
		}
	}

	if options.DryRun || !result.IsChanged() {
		result.Print(true)
		return nil
	}

	total := 0
	uploaded := 0
	skipped := 0

	session := repository.NewSession()
	tasks := []executor.TaskFunc{}
	mtx := sync.Mutex{}

	mapUploadBlob := map[string]DiffRecord{}
	for _, record := range result.Records {
		if record.Type == DiffTypeAdd || record.Type == DiffTypeChange {
			if record.Hash == "" {
				// it is symbolic link. no upload required.
				continue
			}

			if _, ok := mapUploadBlob[record.Hash]; !ok {
				mapUploadBlob[record.Hash] = record
				total++
			}
		}
	}

	for hash, record := range mapUploadBlob {
		h := hash
		p := record.Path
		s := record.Size

		meter := session.NewMeter()

		task := func(ctx context.Context) error {
			uploadResult, err := mngr.UploadBlob(p, h, meter, checkSkip)
			if err != nil {
				return err
			}

			meter.SetBytes(s)
			mtx.Lock()
			uploaded++
			if uploadResult.Skip {
				skipped++
			}
			mtx.Unlock()
			return nil
		}
		tasks = append(tasks, task)
	}

	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()

	done := make(chan error)
	go func() {
		err = executor.ExecuteAll(0, tasks...)
		done <- err
	}()

	stop := false
	for !stop {
		select {
		case err = <-done:
			stop = true
		case <-ticker.C:
		}
		fmt.Printf("upload objects: (%d/%d), skipped: %d, speed: %5v/s    \r", uploaded, total, skipped, session.CalculateSpeed())
	}

	if err != nil {
		return err
	}

	fmt.Println()
	result.Print(false)

	_, hash := MakeCommitMetadata(commit)
	fmt.Println("create commit: " + hash)
	err = mngr.Commit(*commit)
	if err != nil {
		return err
	}

	fmt.Println("update ref: latest -> " + hash)
	err = mngr.AddRef(RefLatest, hash)
	if err != nil {
		return err
	}

	if options.Tag != nil {
		tag := *options.Tag
		err = mngr.AddTag(hash, tag)
		if err != nil {
			return err
		}

		fmt.Println("add tag: " + tag + " -> " + hash)
	}

	return nil
}

func (mngr *ArtifactManager) MakeEmptyCommit() *Commit {
	return &Commit{
		CreatedAt: time.Now(),
		Parent:    "",
		Message:   nil,
		Blobs:     []BlobMetaData{},
	}
}

func (mngr *ArtifactManager) MakeWorkspaceCommit(parent string, message *string, filter func(path string) bool) (*Commit, error) {
	baseDir := mngr.baseDir
	commit := Commit{
		CreatedAt: time.Now(),
		Parent:    parent,
		Message:   message,
		Blobs:     make([]BlobMetaData, 0),
	}

	tasks := []executor.TaskFunc{}
	mutex := sync.Mutex{}

	err := filepath.Walk(baseDir, func(absPath string, info fs.FileInfo, err error) error {
		if err != nil {
			return ErrWorkspaceNotFound
		}
		if info.IsDir() {
			return nil
		}

		path := absPath[len(baseDir)+1:]
		if strings.HasPrefix(path, ".avc/") {
			return nil
		}

		if filter != nil && !filter(path) {
			return nil
		}

		if info.Mode()&os.ModeSymlink != 0 {
			// symbolic
		} else if !info.Mode().IsRegular() {
			return fmt.Errorf("not supported file type. %s -> %x", path, info.Mode().Type())
		}

		task := func(ctx context.Context) error {
			metadata, err := MakeBlobMetadata(baseDir, path)
			if err != nil {
				return fmt.Errorf("cannot make metadata: %s", path)
			}

			mutex.Lock()
			commit.Blobs = append(commit.Blobs, metadata)
			mutex.Unlock()

			return nil
		}
		tasks = append(tasks, task)
		return nil
	})
	if err != nil {
		return nil, err
	}

	err = executor.ExecuteAll(0, tasks...)
	if err != nil {
		return nil, err
	}

	return &commit, nil
}

func (mngr *ArtifactManager) Pull(options PullOptions) error {
	var err error
	if !options.NoFetch {
		err = mngr.Fetch()
		if err != nil {
			return err
		}
	}

	refOrCommit := RefLatest
	if options.RefOrCommit != nil {
		refOrCommit = *options.RefOrCommit
	}

	// Make the remote commit
	log.Debugln("get the remote commit")
	commitHash, err := mngr.FindCommitOrReference(refOrCommit)
	if err != nil {
		var refPath string
		if refOrCommit == RefLatest {
			refPath = RefLatest
		} else {
			refPath = "tags/" + refOrCommit
		}
		commitHash, err = mngr.GetRef(refPath)
		if err != nil {
			if refOrCommit == RefLatest {
				return ErrEmptyRepository
			} else {
				return ReferenceNotFoundError{Ref: refOrCommit}
			}
		}
	}

	commitRemote, err := mngr.GetCommit(commitHash)
	if err != nil && err != ErrEmptyRepository {
		return err
	}

	// Get the local commit hash
	log.Debugln("make the local commit")
	avcIgnore, err := NewAvcIgnore(mngr.baseDir)
	avcIgnoreFilter := func(path string) bool {
		return true
	}

	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		avcIgnoreFilter = func(path string) bool {
			return !avcIgnore.MatchesPath(path)
		}
	}
	commitLocal, err := mngr.MakeWorkspaceCommit("", nil, avcIgnoreFilter)
	if err != nil {
		if err != ErrWorkspaceNotFound {
			return err
		} else {
			commitLocal = mngr.MakeEmptyCommit()
		}
	}

	// Diff
	log.Debugln("diff")
	result, err := mngr.Diff(DiffOptions{
		NoDelete:      !options.Delete,
		LeftCommit:    commitLocal,
		RightCommit:   commitRemote,
		AddFilter:     avcIgnoreFilter,
		ChangeFilter:  avcIgnoreFilter,
		DeleteFilter:  avcIgnoreFilter,
		IncludeFilter: options.FileFilter,
	})
	if err != nil {
		return err
	}

	if options.DryRun || !result.IsChanged() {
		result.Print(true)
		return nil
	}

	// download
	log.Debugln("download")
	total := 0
	downloaded := 0

	session := repository.NewSession()
	tasks := []executor.TaskFunc{}
	mtx := sync.Mutex{}
	for _, record := range result.Records {
		if record.Type != DiffTypeAdd && record.Type != DiffTypeChange {
			continue
		}

		if record.Hash == "" {
			//symbolic link
			continue
		}

		if record.Type == DiffTypeChange && record.OldHash == record.Hash {
			// mode change
			continue
		}

		p := record.Path
		h := record.Hash
		s := record.Size

		task := func(ctx context.Context) error {
			meter := session.NewMeter()
			_, err := mngr.DownloadBlob(p, h, meter)
			if err != nil {
				return err
			}
			mtx.Lock()
			downloaded++
			mtx.Unlock()
			meter.SetBytes(s)

			return nil
		}

		tasks = append(tasks, task)
		total++
	}

	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()

	done := make(chan error)
	go func() {
		err = executor.ExecuteAll(10, tasks...)
		done <- err
	}()

	stop := false
	for !stop {
		select {
		case err = <-done:
			stop = true
		case <-ticker.C:
		}
		fmt.Printf("download objects: (%d/%d), speed: %5v/s    \r", downloaded, total, session.CalculateSpeed())
	}
	fmt.Println()

	// delete, rename, symlink, chmod
	log.Debugln("delete, rename, symlink, chmod")
	for _, record := range result.Records {
		absPath := filepath.Join(mngr.baseDir, record.Path)
		mode := record.Mode

		switch record.Type {
		case DiffTypeAdd:
			if record.Link != "" {
				err := symlinkFile(record.Link, absPath)
				if err != nil {
					return err
				}
			} else {
				err := chmod(absPath, mode)
				if err != nil {
					return err
				}
			}
		case DiffTypeChange:
			if record.Link != "" {
				err := deleteFile(absPath)
				if err != nil {
					return err
				}

				err = symlinkFile(record.Link, absPath)
				if err != nil {
					return err
				}
			} else {
				err := chmod(absPath, mode)
				if err != nil {
					return err
				}
			}
		case DiffTypeDelete:
			err := deleteFile(absPath)
			if err != nil {
				return err
			}
		case DiffTypeRename:
			err := renameFile(filepath.Join(mngr.baseDir, record.OldPath), absPath)
			if err != nil {
				return err
			}

			if record.Hash != "" {
				err := chmod(absPath, mode)
				if err != nil {
					return err
				}
			}
		}
	}

	if options.Delete {
		_, err = removeEmptyDirs(mngr.baseDir, false)
		if err != nil {
			return err
		}
	}
	_, err = removeEmptyDirs(filepath.Join(mngr.baseDir, ".avc"), true)
	if err != nil {
		return err
	}

	// print summary
	result.Print(false)

	return nil
}

func (mngr *ArtifactManager) ListTags() error {
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

func (mngr *ArtifactManager) AddTag(refOrCommit, tag string) error {
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

func (mngr *ArtifactManager) DeleteTag(tag string) error {
	if tag == RefLatest {
		return errors.New("latest cannot be a tag")
	}

	err := mngr.DeleteRef("tags/" + tag)
	if err != nil {
		return err
	}

	return nil
}

func (mngr *ArtifactManager) List(refOrCommit string) error {
	err := mngr.Fetch()
	if err != nil {
		return err
	}

	commitHash, err := mngr.FindCommitOrReference(refOrCommit)
	if err != nil {
		return err
	}

	commit, err := mngr.GetCommit(commitHash)
	if err != nil {
		return err
	}

	sort.Slice(commit.Blobs, func(i, j int) bool {
		return commit.Blobs[i].Path < commit.Blobs[j].Path
	})

	for _, blob := range commit.Blobs {
		fmt.Println(blob.Path)
	}

	return nil
}

func (mngr *ArtifactManager) Diff(option DiffOptions) (DiffResult, error) {
	type DiffEntry struct {
		left  *BlobMetaData
		right *BlobMetaData
	}
	entries := map[string]DiffEntry{}

	var commitHash string
	var err error

	// Step 1: Prepare the left and right commits

	// left
	leftCommit := option.LeftCommit
	if leftCommit == nil {
		commitHash, err = mngr.FindCommitOrReference(option.LeftRef)
		if err != nil {
			return DiffResult{}, err
		}

		leftCommit, err = mngr.GetCommit(commitHash)
		if err != nil {
			return DiffResult{}, err
		}
	}
	for i, blob := range leftCommit.Blobs {
		if option.IncludeFilter != nil && !option.IncludeFilter(blob.Path) {
			continue
		}
		entry := entries[blob.Path]
		entry.left = &leftCommit.Blobs[i]
		entries[blob.Path] = entry
	}

	// right
	rightCommit := option.RightCommit
	if rightCommit == nil {
		commitHash, err = mngr.FindCommitOrReference(option.RightRef)
		if err != nil {
			return DiffResult{}, err
		}

		rightCommit, err = mngr.GetCommit(commitHash)
		if err != nil {
			return DiffResult{}, err
		}
	}

	for i, blob := range rightCommit.Blobs {
		if option.IncludeFilter != nil && !option.IncludeFilter(blob.Path) {
			continue
		}
		entry := entries[blob.Path]
		entry.right = &rightCommit.Blobs[i]
		entries[blob.Path] = entry
	}

	// Step 2: Compare left and right and create the changesets (added, deleted, changed)
	mapAdded := map[string][]DiffRecord{}
	mapDeleted := map[string][]DiffRecord{}
	mapChanged := map[string][]DiffRecord{}

	appendOrMake := func(s []DiffRecord, item DiffRecord) []DiffRecord {
		if s != nil {
			return append(s, item)
		} else {
			return []DiffRecord{item}
		}
	}

	for path, entry := range entries {
		if entry.left == nil && entry.right != nil {
			if option.AddFilter != nil {
				if !option.AddFilter(path) {
					continue
				}
			}

			var key string
			record := DiffRecord{Type: DiffTypeAdd, Path: path}
			if entry.right.Link != "" {
				record.Link = entry.right.Link
				key = record.Link
			} else {
				record.Hash = entry.right.Hash
				record.Size = entry.right.Size
				record.Mode = entry.right.Mode
				key = record.Hash
			}

			mapAdded[key] = appendOrMake(mapAdded[key], record)
		} else if entry.left != nil && entry.right == nil {
			if option.NoDelete {
				continue
			}

			if option.DeleteFilter != nil {
				if !option.DeleteFilter(path) {
					continue
				}
			}

			var key string
			record := DiffRecord{Type: DiffTypeDelete, Path: path}
			if entry.left.Link != "" {
				record.Link = entry.left.Link
				key = record.Link
			} else {
				record.Hash = entry.left.Hash
				record.Size = entry.left.Size
				record.Mode = entry.left.Mode
				key = record.Hash
			}
			mapDeleted[key] = appendOrMake(mapDeleted[key], record)
		} else if entry.left.Hash != entry.right.Hash || entry.left.Mode != entry.right.Mode || entry.left.Link != entry.right.Link {
			if option.ChangeFilter != nil {
				if !option.ChangeFilter(path) {
					continue
				}
			}

			var key string
			record := DiffRecord{
				Type: DiffTypeChange,
				Path: path,
			}

			if entry.left.Link != "" {
				record.OldLink = entry.left.Link
				key = record.OldLink
			} else {
				record.OldHash = entry.left.Hash
				record.OldSize = entry.left.Size
				record.OldMode = entry.left.Mode
				key = record.OldHash
			}

			if entry.right.Link != "" {
				record.Link = entry.right.Link
			} else {
				record.Hash = entry.right.Hash
				record.Size = entry.right.Size
				record.Mode = entry.right.Mode
			}

			mapChanged[key] = appendOrMake(mapChanged[key], record)
		}
	}

	// Step3: Merge "added" and "deleted" with the same content hash to "rename"
	mapRenamed := map[string][]DiffRecord{}
	for key, addedPaths := range mapAdded {
		deleledPaths := mapDeleted[key]
		if deleledPaths == nil {
			continue
		}

		// n = min(added, deleted)
		n := len(addedPaths)
		if len(deleledPaths) < n {
			n = len(deleledPaths)
		}

		//
		renamedRecords := []DiffRecord{}
		for i := 0; i < n; i++ {

			record := DiffRecord{
				Type: DiffTypeRename,
				Path: addedPaths[i].Path,
				Hash: addedPaths[i].Hash,
				Link: addedPaths[i].Link,
				Mode: addedPaths[i].Mode,
				Size: addedPaths[i].Size,

				OldPath: deleledPaths[i].Path,
				OldHash: deleledPaths[i].Hash,
				OldLink: deleledPaths[i].Link,
				OldMode: deleledPaths[i].Mode,
				OldSize: deleledPaths[i].Size,
			}

			renamedRecords = append(renamedRecords, record)
		}
		mapAdded[key] = addedPaths[n:]
		mapDeleted[key] = deleledPaths[n:]
		mapRenamed[key] = renamedRecords
	}

	// Step 4: Merge the the 4 maps to the diff record list
	records := []DiffRecord{}
	for _, added := range mapAdded {
		records = append(records, added...)
	}

	for _, deleted := range mapDeleted {
		records = append(records, deleted...)
	}

	for _, changed := range mapChanged {
		records = append(records, changed...)
	}

	for _, renamed := range mapRenamed {
		records = append(records, renamed...)
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].Path < records[j].Path
	})

	return DiffResult{
		Records: records,
	}, nil
}

func (mngr *ArtifactManager) Status() (DiffResult, error) {

	// Make the remote latest commit
	commitHash, err := mngr.FindCommitOrReference(RefLatest)
	if err != nil {
		commitHash, err = mngr.GetRef(RefLatest)
		if err != nil {
			return DiffResult{}, ErrEmptyRepository
		}
	}
	commitRemote, err := mngr.GetCommit(commitHash)
	if err != nil && err != ErrEmptyRepository {
		return DiffResult{}, err
	}

	// Get the local commit hash
	avcIgnore, err := NewAvcIgnore(mngr.baseDir)
	avcIgnoreFilter := func(path string) bool {
		return true
	}

	if err != nil {
		if !os.IsNotExist(err) {
			return DiffResult{}, err
		}
	} else {
		avcIgnoreFilter = func(path string) bool {
			return !avcIgnore.MatchesPath(path)
		}
	}

	commitLocal, err := mngr.MakeWorkspaceCommit("", nil, avcIgnoreFilter)
	if err != nil {
		if err != ErrWorkspaceNotFound {
			return DiffResult{}, err
		} else {
			commitLocal = mngr.MakeEmptyCommit()
		}
	}

	// Diff
	result, err := mngr.Diff(DiffOptions{
		LeftCommit:   commitRemote,
		RightCommit:  commitLocal,
		AddFilter:    avcIgnoreFilter,
		ChangeFilter: avcIgnoreFilter,
		DeleteFilter: avcIgnoreFilter,
	})
	if err != nil {
		return DiffResult{}, err
	}

	return result, nil
}

func (mngr *ArtifactManager) Log(refOrCommit string) error {
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
		return ErrEmptyRepository
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

func (result DiffResult) IsChanged() bool {
	return len(result.Records) > 0
}

func (result DiffResult) IsAppendOnly() bool {
	var modified int
	records := result.Records
	if records == nil {
		records = []DiffRecord{}
	}

	for _, record := range records {
		if record.Type != DiffTypeAdd {
			modified++
		}
	}

	return modified == 0
}

func (result *DiffResult) Print(verbose bool) {
	records := result.Records
	if records == nil {
		records = []DiffRecord{}
	}

	var added, deleted, modified, renamed int

	for _, record := range records {
		switch record.Type {
		case DiffTypeAdd:
			if verbose {
				color.HiGreen(fmt.Sprintf("+ %s\n", record.Path))
			}
			added++
		case DiffTypeDelete:
			if verbose {
				color.HiRed(fmt.Sprintf("- %s\n", record.Path))
			}
			deleted++
		case DiffTypeChange:
			if verbose {
				color.HiYellow(fmt.Sprintf("M %s\n", record.Path))
			}
			modified++
		case DiffTypeRename:
			if verbose {
				color.HiYellow(fmt.Sprintf("R %s -> %s\n", record.OldPath, record.Path))
			}
			renamed++
		}
	}

	if !result.IsChanged() {
		fmt.Println("no changed")
	} else {
		fmt.Printf("%d modified(M), %d added(+), %d deleted(-), %d renamed(R)\n", modified, added, deleted, renamed)
	}
}
