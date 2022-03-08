package core

import (
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

	_, err = os.Stat(filepath.Join(wp2, ".art"))
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

	_, err := os.Stat(filepath.Join(wp2, ".art/config"))
	assert.False(t, os.IsNotExist(err))
}

func TestPushWithIgnore(t *testing.T) {
	wp1 := t.TempDir()
	wp2 := t.TempDir()
	repo := t.TempDir()

	writeFile([]byte("a"), filepath.Join(wp1, "a"))
	writeFile([]byte("b"), filepath.Join(wp1, "b"))
	writeFile([]byte("c"), filepath.Join(wp1, "c"))

	artIgnore := `
^a$
^e$
`

	writeFile([]byte(artIgnore), filepath.Join(wp1, ".artignore"))

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
	artIgnore := `
^a$
^e$
`
	writeFile([]byte(artIgnore), filepath.Join(wp2, ".artignore"))
	writeFile([]byte("abc"), filepath.Join(wp2, "a"))
	writeFile([]byte("efg"), filepath.Join(wp2, "e"))
	InitWorkspace(wp2, repo)
	config, _ = LoadConfig(wp2)
	mngr2, _ := NewArtifactManager(config)
	err = mngr2.Pull(PullOptions{Mode: ChangeModeMerge})
	assert.Empty(t, err)

	data, _ := readFile(filepath.Join(wp2, "a"))
	assert.Equal(t, "abc", string(data))
	data, _ = readFile(filepath.Join(wp2, "c"))
	assert.Equal(t, "c", string(data))
	data, _ = readFile(filepath.Join(wp2, "e"))
	assert.Equal(t, "efg", string(data))
}
