package util

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func Sha256SumFromFile(path string) string {
	hasher := sha256.New()
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("error")
	}
	defer f.Close()
	if _, err := io.Copy(hasher, f); err != nil {
		fmt.Println("error")
	}
	sum := hasher.Sum([]byte(""))
	return fmt.Sprintf("%x", sum)
}

func Sha1SumFromFile(path string) string {
	hasher := sha1.New()
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("error")
	}
	defer f.Close()
	if _, err := io.Copy(hasher, f); err != nil {
		fmt.Println("error")
	}
	sum := hasher.Sum([]byte(""))
	return fmt.Sprintf("%x", sum)
}
