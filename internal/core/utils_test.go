package core

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveEmptyDirs(t *testing.T) {
	tmpDir := t.TempDir()

	assert.NoError(t, os.Mkdir(tmpDir+"/a", os.ModePerm))
	removed, err := removeEmptyDirs(tmpDir+"/a", true)
	assert.True(t, removed)
	assert.Empty(t, err)

	assert.NoError(t, os.Mkdir(tmpDir+"/b", os.ModePerm))
	assert.NoError(t, writeFile([]byte("hello"), tmpDir+"/b/hello"))
	removed, err = removeEmptyDirs(tmpDir+"/b", true)
	assert.False(t, removed)
	assert.Empty(t, err)

	assert.NoError(t, os.MkdirAll(tmpDir+"/c/c/c/c/a", os.ModePerm))
	assert.NoError(t, os.MkdirAll(tmpDir+"/c/c/c/c/b", os.ModePerm))
	assert.NoError(t, os.MkdirAll(tmpDir+"/c/c/c/c/c", os.ModePerm))
	assert.NoError(t, writeFile([]byte("hello"), tmpDir+"/c/c/a"))
	removed, err = removeEmptyDirs(tmpDir+"/c", true)
	assert.Empty(t, err)
	assert.False(t, removed)
	stat, err := os.Stat(tmpDir + "/c/c/a")
	assert.Equal(t, "a", stat.Name())
	assert.Empty(t, err)
	_, err = os.Stat("/c/c/c")
	assert.Error(t, err)
}
