package core

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPutGet(t *testing.T) {
	tempDir := t.TempDir()
	wp1 := tempDir + "/wp1"
	meta1 := tempDir + "/meta1"
	wp2 := tempDir + "/wp2"
	meta2 := tempDir + "/meta2"
	repo := tempDir + "/repo"

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
	t.TempDir()
	wp1 := t.TempDir() + "/wp1"
	wp2 := t.TempDir() + "/wp2"
	repo := t.TempDir() + "/repo"

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
