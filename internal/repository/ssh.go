package repository

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Local Filesystem
type SSHRepository struct {
	Host    string
	BaseDir string
}

func NewSSHRepository(host, basePath string) (*SSHRepository, error) {
	cmd := exec.Command("ssh", "-V")
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	return &SSHRepository{
		Host:    host,
		BaseDir: basePath,
	}, nil
}

func (repo *SSHRepository) Upload(localPath, repoPath string, m *Meter) error {
	path := filepath.Join(repo.BaseDir, repoPath)
	script := `
set -e
mkdir -p ${DEST_DIR}
cat | gzip -cd > ${DEST_PATH}
`
	expandMap := map[string]string{
		"DEST_DIR":  filepath.Dir(path),
		"DEST_PATH": path,
	}

	script = os.Expand(script, func(k string) string {
		return expandMap[k]
	})

	cmd := repo.rcommand(script)
	dest, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	var stderr = bytes.Buffer{}
	cmd.Stderr = &stderr

	source, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer source.Close()

	gdest := gzip.NewWriter(dest)
	defer gdest.Close()

	go func() {
		CopyWithMeter(gdest, source, m)
		source.Close()
		gdest.Close()
		dest.Close()
	}()

	err = cmd.Run()
	if err != nil {
		return errors.New(stderr.String())
	}

	return nil
}

func (repo *SSHRepository) Download(repoPath, localPath string, m *Meter) error {
	path := filepath.Join(repo.BaseDir, repoPath)
	cmd := repo.rcommand(fmt.Sprintf("gzip -c %s", path))

	dest, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	src, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	go func() {
		gsrc, err := gzip.NewReader(src)
		if err != nil {
			return
		}
		defer gsrc.Close()

		CopyWithMeter(dest, gsrc, m)
		dest.Close()
		gsrc.Close()
		src.Close()
	}()

	err = cmd.Run()
	if err != nil {
		return errors.New(stderr.String())
	}

	return nil
}

func (repo *SSHRepository) Delete(repoPath string) error {
	path := filepath.Join(repo.BaseDir, repoPath)
	cmd := repo.rcommand("rm " + path)
	_, err := cmd.Output()
	if err != nil {
		return err
	}

	return nil
}

func (repo *SSHRepository) Stat(repoPath string) (FileInfo, error) {
	path := filepath.Join(repo.BaseDir, repoPath)
	cmd := repo.rcommand("ls -ald " + path)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	result, err := sshParseListRecord(string(out))
	if err != nil {
		return nil, err
	}
	result.name = filepath.Base(repoPath)

	return result, nil
}

func (repo *SSHRepository) List(repoPath string) ([]FileInfo, error) {
	entries := make([]FileInfo, 0)
	path := filepath.Join(repo.BaseDir, repoPath)
	cmd := repo.rcommand("ls -al " + path)

	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Run()
	if err != nil {
		return entries, nil
	}

	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "total") {
			continue
		}

		info, err := sshParseListRecord(scanner.Text())
		if err != nil {
			return nil, err
		}

		if info.Name() == "." || info.Name() == ".." {
			continue
		}
		entries = append(entries, info)
	}

	return entries, nil
}

func (repo *SSHRepository) rcommand(script string) *exec.Cmd {
	return exec.Command("ssh", repo.Host, script)
}

func sshParseListRecord(record string) (SSHFileInfo, error) {
	components := strings.Split(record, " ")
	if len(components) < 7 {
		return SSHFileInfo{}, os.ErrInvalid
	}

	mode := components[0]
	name := components[len(components)-1]
	return SSHFileInfo{
		name:  name,
		isDir: mode[0] == 'd',
	}, nil
}

type SSHFileInfo struct {
	name  string
	isDir bool
}

func (e SSHFileInfo) Name() string {
	return e.name
}

func (e SSHFileInfo) IsDir() bool {
	return e.isDir
}
