package core

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPutGet(t *testing.T) {
	wp1 := t.TempDir()
	meta1 := t.TempDir()
	wp2 := t.TempDir()
	meta2 := t.TempDir()
	repo := t.TempDir()

	path := "test"
	content := "test-data"

	writeFile([]byte(content), filepath.Join(wp1, path))

	config := NewConfig(wp1, meta1, repo)
	mngr1, err := NewArtifactManager(config)
	assert.Empty(t, err)
	err = mngr1.Push(PushOptions{})
	assert.Empty(t, err)

	config = NewConfig(wp2, meta2, repo)
	mngr2, err := NewArtifactManager(config)
	assert.Empty(t, err)
	err = mngr2.Pull(PullOptions{})
	assert.Empty(t, err)

	data, err := readFile(filepath.Join(wp2, path))
	assert.Empty(t, err)
	assert.Equal(t, string(data), content)

	_, err = os.Stat(filepath.Join(wp2, ".avc"))
	assert.True(t, os.IsNotExist(err))
}

func TestPushPull(t *testing.T) {
	wp1 := t.TempDir()
	wp2 := t.TempDir()
	repo := t.TempDir()

	path := "test"
	content := "test-data"

	writeFile([]byte(content), filepath.Join(wp1, path))

	InitWorkspace(wp1, repo)
	config, _ := LoadConfig(wp1)
	mngr1, _ := NewArtifactManager(config)
	mngr1.Push(PushOptions{})

	InitWorkspace(wp2, repo)
	config, _ = LoadConfig(wp2)
	mngr2, _ := NewArtifactManager(config)
	mngr2.Pull(PullOptions{})

	data, _ := readFile(filepath.Join(wp2, path))
	assert.Equal(t, string(data), content)

	_, err := os.Stat(filepath.Join(wp2, ".avc/config"))
	assert.False(t, os.IsNotExist(err))
}

func TestPushWithIgnore(t *testing.T) {
	wp1 := t.TempDir()
	wp2 := t.TempDir()
	repo := t.TempDir()

	writeFile([]byte("a"), filepath.Join(wp1, "a"))
	writeFile([]byte("b"), filepath.Join(wp1, "b"))
	writeFile([]byte("c"), filepath.Join(wp1, "c"))

	avcIgnore := `
a
e
`

	writeFile([]byte(avcIgnore), filepath.Join(wp1, ".avcignore"))

	InitWorkspace(wp1, repo)
	config, _ := LoadConfig(wp1)
	mngr1, _ := NewArtifactManager(config)
	err := mngr1.Push(PushOptions{})
	assert.Empty(t, err)

	InitWorkspace(wp2, repo)
	config, _ = LoadConfig(wp2)
	mngr2, _ := NewArtifactManager(config)
	err = mngr2.Pull(PullOptions{})
	assert.Empty(t, err)

	data, _ := readFile(filepath.Join(wp2, "a"))
	assert.Equal(t, "", string(data))
	data, _ = readFile(filepath.Join(wp2, "b"))
	assert.Equal(t, "b", string(data))
	data, _ = readFile(filepath.Join(wp2, "c"))
	assert.Equal(t, "c", string(data))
}

func TestPullWithIgnore(t *testing.T) {
	wp1 := t.TempDir()
	wp2 := t.TempDir()
	repo := t.TempDir()

	// push
	writeFile([]byte("a"), filepath.Join(wp1, "a"))
	writeFile([]byte("b"), filepath.Join(wp1, "b"))
	writeFile([]byte("c"), filepath.Join(wp1, "c"))
	InitWorkspace(wp1, repo)
	config, _ := LoadConfig(wp1)
	mngr1, _ := NewArtifactManager(config)
	err := mngr1.Push(PushOptions{})
	assert.Empty(t, err)

	// pull
	avcIgnore := `
a
e
`
	writeFile([]byte(avcIgnore), filepath.Join(wp2, ".avcignore"))
	writeFile([]byte("abc"), filepath.Join(wp2, "a"))
	writeFile([]byte("efg"), filepath.Join(wp2, "e"))
	InitWorkspace(wp2, repo)
	config, _ = LoadConfig(wp2)
	mngr2, _ := NewArtifactManager(config)
	err = mngr2.Pull(PullOptions{})
	assert.Empty(t, err)

	data, _ := readFile(filepath.Join(wp2, "a"))
	assert.Equal(t, "abc", string(data))
	data, _ = readFile(filepath.Join(wp2, "c"))
	assert.Equal(t, "c", string(data))
	data, _ = readFile(filepath.Join(wp2, "e"))
	assert.Equal(t, "efg", string(data))
}

func TestSymlink(t *testing.T) {
	filepath.Walk("/tmp/symlink/", func(absPath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Mode()&fs.ModeSymlink > 0 {
			link, err := os.Readlink(absPath)
			if err != nil {
				return err
			}
			fmt.Printf("symlink: %o %s -> %s\n", info.Mode(), info.Name(), link)
		} else if info.IsDir() {
			fmt.Printf("dir: %o %s\n", info.Mode(), info.Name())
		} else {
			fmt.Printf("file: %o %s\n", info.Mode(), info.Name())
		}

		return nil
	})

	os.Symlink("/tmp/symlink/target", "/tmp/symlink/link")
}
