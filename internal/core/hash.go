package core

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func Sha256Sum(content []byte) string {
	sum := sha256.New().Sum(content)
	return fmt.Sprintf("%x", sum)
}

func Sha256SumFromFile(path string) (string, error) {
	hasher := sha256.New()
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}
	sum := hasher.Sum([]byte{})
	return fmt.Sprintf("%x", sum), nil
}

func Sha1Sum(content []byte) string {
	hasher := sha1.New()
	hasher.Write(content)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func Sha1SumFromFile(path string) (string, error) {
	hasher := sha1.New()
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}
	sum := hasher.Sum([]byte{})
	return fmt.Sprintf("%x", sum), nil
}
