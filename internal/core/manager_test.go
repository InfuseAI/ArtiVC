package core

// func TestUpload(t *testing.T) {
// 	baseDir := t.TempDir() + "/repo"
// 	repoDir := t.TempDir() + "/data"
// 	path := "test"
// 	content := "test-data"
// 	repo := repository.NewLocalFileSystemRepository() LocalFileSystemRepository{
// 		BaseDir: baseDir,
// 		RepoDir: repoDir,
// 	}

// 	// prepare data
// 	fullPath := filepath.Join(repo.BaseDir, path)
// 	os.MkdirAll(filepath.Dir(fullPath), fs.ModePerm)
// 	if err := ioutil.WriteFile(fullPath, []byte(content), 0o644); err != nil {
// 		assert.Fail(t, "cannot write file")
// 		return
// 	}

// 	// upload
// 	// get the size
// 	metaData, err := core.MakeBlobMetadata(baseDir, path)
// 	if err != nil {
// 		assert.Fail(t, "cannot create metedata")
// 		return
// 	}

// 	err = repo.UploadBlob(metaData)
// 	if err != nil {
// 		assert.Fail(t, "cannot create metedata")
// 		return
// 	}

// 	destFile := GenObjectPath(repo.RepoDir, metaData.Hash)
// 	data, _ := ioutil.ReadFile(destFile)
// 	assert.Equal(t, content, string(data))
// }
