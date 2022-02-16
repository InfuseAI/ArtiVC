package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPutGet(t *testing.T) {
	t.TempDir()
	wp1 := t.TempDir() + "/wp1"
	meta1 := t.TempDir() + "/meta1"
	wp2 := t.TempDir() + "/wp2"
	meta2 := t.TempDir() + "/meta2"
	repo := t.TempDir() + "/repo"

	path := "test"
	content := "test-data"

	writeFile([]byte(content), filepath.Join(wp1, path))

	config := NewConfig(wp1, meta1, repo)
	mngr1, _ := NewArtifactManager(config)
	mngr1.Push()

	config = NewConfig(wp2, meta2, repo)
	mngr2, _ := NewArtifactManager(config)
	mngr2.Pull()

	data, _ := readFile(filepath.Join(wp2, path))
	assert.Equal(t, string(data), content)

	_, err := os.Stat(filepath.Join(wp2, ".art"))
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
	mngr1.Push()

	InitWorkspace(wp2, repo)
	config, _ = LoadConfig(wp2)
	mngr2, _ := NewArtifactManager(config)
	mngr2.Pull()

	data, _ := readFile(filepath.Join(wp2, path))
	assert.Equal(t, string(data), content)

	_, err := os.Stat(filepath.Join(wp2, ".art/config"))
	assert.False(t, os.IsNotExist(err))
}

func TestTemp(t *testing.T) {

	a, _ := ioutil.TempDir("", "example")
	b := os.TempDir()
	fmt.Printf(a)
	fmt.Printf(b)

}
