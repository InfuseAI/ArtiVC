package core

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

func MakeObjectPath(hash string) string {
	return fmt.Sprintf("objects/%s/%s", hash[:2], hash[2:])
}

func MakeCommitPath(hash string) string {
	return fmt.Sprintf("commits/%s", hash)
}

func MakeRefPath(ref string) string {
	return fmt.Sprintf("refs/%s", ref)
}

func MakeTagPath(ref string) string {
	return fmt.Sprintf("refs/tags/%s", ref)
}

func mkdirsForFile(file string) error {
	return os.MkdirAll(filepath.Dir(file), fs.ModePerm)
}

func removeEmptyDirs(dir string, removeSelf bool) (bool, error) {
	var hasEntries bool

	entires, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}
	for _, entry := range entires {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			removed, err := removeEmptyDirs(subdir, true)
			if err != nil {
				return false, err
			}
			if !removed {
				hasEntries = true
			}
		} else {
			hasEntries = true
		}
	}

	if !hasEntries && removeSelf {
		err := os.Remove(dir)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

func writeFile(content []byte, dst string) error {
	err := os.MkdirAll(filepath.Dir(dst), fs.ModePerm)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dst, content, 0o644)
	return err
}

func readFile(src string) ([]byte, error) {
	return ioutil.ReadFile(src)
}

func writeGzipFile(content []byte, dst string) error {
	err := os.MkdirAll(filepath.Dir(dst), fs.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(dst)
	if err != nil {
		return err
	}
	gfile := gzip.NewWriter(file)
	defer gfile.Close()
	_, err = gfile.Write(content)
	return err
}

func readGzipFile(src string) ([]byte, error) {
	file, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	gfile, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer gfile.Close()

	return io.ReadAll(gfile)
}

func deleteFile(src string) error {
	return os.Remove(src)
}

func renameFile(src, dst string) error {
	err := os.MkdirAll(filepath.Dir(dst), fs.ModePerm)
	if err != nil {
		return err
	}

	return os.Rename(src, dst)
}

func symlinkFile(target, src string) error {
	return os.Symlink(target, src)
}
