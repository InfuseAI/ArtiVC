package core

import (
	"fmt"
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

func mkdirsForFile(file string) error {
	return os.MkdirAll(filepath.Dir(file), fs.ModePerm)
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
