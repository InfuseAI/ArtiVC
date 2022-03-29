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
	var err error
	wp1 := t.TempDir()
	wp2 := t.TempDir()
	repo := t.TempDir()

	// Firt version
	// a = "a"
	// b -> bb
	// c -> cc
	// d -> dd
	writeFile([]byte("a"), filepath.Join(wp1, "a"))
	symlinkFile("bb", filepath.Join(wp1, "b"))
	symlinkFile("cc", filepath.Join(wp1, "c"))
	symlinkFile("dd", filepath.Join(wp1, "d"))

	InitWorkspace(wp1, repo)
	config, _ := LoadConfig(wp1)
	mngr1, _ := NewArtifactManager(config)
	mngr1.Push(PushOptions{})

	InitWorkspace(wp2, repo)
	config, _ = LoadConfig(wp2)
	mngr2, _ := NewArtifactManager(config)
	mngr2.Pull(PullOptions{})

	data, _ := readFile(filepath.Join(wp2, "a"))
	assert.Equal(t, "a", string(data))
	link, _ := readlinkFile(filepath.Join(wp2, "b"))
	assert.Equal(t, "bb", link)

	// Second version
	// a = "a"    =>  a -> aa
	// b -> bb    => (deleted)
	// c -> cc    =>  c = "c"
	// d -> dd    =>  e -> dd
	deleteFile(filepath.Join(wp1, "a"))
	symlinkFile("aa", filepath.Join(wp1, "a"))
	deleteFile(filepath.Join(wp1, "b"))
	deleteFile(filepath.Join(wp1, "c"))
	writeFile([]byte("c"), filepath.Join(wp1, "c"))
	deleteFile(filepath.Join(wp1, "d"))
	symlinkFile("dd", filepath.Join(wp1, "e"))
	mngr1.Push(PushOptions{})
	mngr2.Pull(PullOptions{Delete: true})

	link, _ = readlinkFile(filepath.Join(wp2, "a"))
	assert.Equal(t, "aa", link)
	_, err = readlinkFile(filepath.Join(wp2, "b"))
	assert.Error(t, err)
	data, _ = readFile(filepath.Join(wp2, "c"))
	assert.Equal(t, "c", string(data))
	_, err = readlinkFile(filepath.Join(wp2, "d"))
	assert.Error(t, err)
	link, _ = readlinkFile(filepath.Join(wp2, "e"))
	assert.Equal(t, "dd", link)
}

func TestPermMode(t *testing.T) {
	wp1 := t.TempDir()
	wp2 := t.TempDir()
	repo := t.TempDir()

	// Firt version
	// a = "a" 644
	// b = "b" 600
	// c = "c" 755
	writeFile([]byte("a"), filepath.Join(wp1, "a"))
	chmod(filepath.Join(wp1, "a"), 0644)
	writeFile([]byte("b"), filepath.Join(wp1, "b"))
	chmod(filepath.Join(wp1, "b"), 0600)
	writeFile([]byte("c"), filepath.Join(wp1, "c"))
	chmod(filepath.Join(wp1, "c"), 0755)

	InitWorkspace(wp1, repo)
	config, _ := LoadConfig(wp1)
	mngr1, _ := NewArtifactManager(config)
	mngr1.Push(PushOptions{})

	InitWorkspace(wp2, repo)
	config, _ = LoadConfig(wp2)
	mngr2, _ := NewArtifactManager(config)
	mngr2.Pull(PullOptions{})

	mode, _ := readFileMode(filepath.Join(wp2, "a"))
	assert.Equal(t, 0644, int(mode))
	mode, _ = readFileMode(filepath.Join(wp2, "b"))
	assert.Equal(t, 0600, int(mode))
	mode, _ = readFileMode(filepath.Join(wp2, "c"))
	assert.Equal(t, 0755, int(mode))

	// Second versio n
	// a = "a" 644 => a  = "a" 755
	// b = "b" 600 => bb = "b" 600
	// c = "c" 755 => cc = "c" 700
	//          (new) d  = "d" 755
	chmod(filepath.Join(wp1, "a"), 0755)
	renameFile(filepath.Join(wp1, "b"), filepath.Join(wp1, "bb"))
	renameFile(filepath.Join(wp1, "c"), filepath.Join(wp1, "cc"))
	chmod(filepath.Join(wp1, "cc"), 0700)
	writeFile([]byte("d"), filepath.Join(wp1, "d"))
	chmod(filepath.Join(wp1, "d"), 0755)

	mngr1.Push(PushOptions{})
	mngr2.Pull(PullOptions{Delete: true})

	mode, _ = readFileMode(filepath.Join(wp2, "a"))
	assert.Equal(t, 0755, int(mode))
	mode, _ = readFileMode(filepath.Join(wp2, "bb"))
	assert.Equal(t, 0600, int(mode))
	mode, _ = readFileMode(filepath.Join(wp2, "cc"))
	assert.Equal(t, 0700, int(mode))
	mode, _ = readFileMode(filepath.Join(wp2, "d"))
	assert.Equal(t, 0755, int(mode))
}
