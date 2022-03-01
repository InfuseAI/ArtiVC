package core

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestRemoveEmptyDirs(t *testing.T) {
	tmpDir := t.TempDir()

	os.Mkdir(tmpDir+"/a", os.ModePerm)
	removed, err := removeEmptyDirs(tmpDir+"/a", true)
	assert.True(t, removed)
	assert.Empty(t, err)

	os.Mkdir(tmpDir+"/b", os.ModePerm)
	writeFile([]byte("hello"), tmpDir+"/b/hello")
	removed, err = removeEmptyDirs(tmpDir+"/b", true)
	assert.False(t, removed)
	assert.Empty(t, err)

	os.MkdirAll(tmpDir+"/c/c/c/c/a", os.ModePerm)
	os.MkdirAll(tmpDir+"/c/c/c/c/b", os.ModePerm)
	os.MkdirAll(tmpDir+"/c/c/c/c/c", os.ModePerm)
	writeFile([]byte("hello"), tmpDir+"/c/c/a")
	removed, err = removeEmptyDirs(tmpDir+"/c", true)
	assert.Empty(t, err)
	assert.False(t, removed)
	stat, err := os.Stat(tmpDir + "/c/c/a")
	assert.Equal(t, "a", stat.Name())
	assert.Empty(t, err)
	_, err = os.Stat("/c/c/c")
	assert.Error(t, err)
}
