package repository

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/kevinburke/ssh_config"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// Local Filesystem
type SSHRepository struct {
	Host       string
	BaseDir    string
	SSHClient  *ssh.Client
	SFTPClient *sftp.Client
}

func NewSSHRepository(host, basePath string) (*SSHRepository, error) {
	authMethods := []ssh.AuthMethod{}

	agentSock := os.Getenv("SSH_AUTH_SOCK")
	if agentSock != "" {
		agentConn, err := net.Dial("unix", agentSock)
		if err != nil {
			return nil, err
		}

		agentClient := agent.NewClient(agentConn)
		authMethods = append(authMethods, ssh.PublicKeysCallback(agentClient.Signers))
	}

	// load ssh config
	f, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "config"))
	if err != nil {
		return nil, err
	}
	cfg, err := ssh_config.Decode(f)
	if err != nil {
		return nil, err
	}

	hostname, _ := cfg.Get(host, "Hostname")
	port, _ := cfg.Get(host, "Port")
	user, _ := cfg.Get(host, "User")

	config := &ssh.ClientConfig{
		User:            user,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// setup the ssh client and sftp client
	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", hostname, port), config)
	if err != nil {
		return nil, err
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, err
	}

	rand.Seed(time.Now().UnixNano())

	return &SSHRepository{
		Host:       host,
		BaseDir:    basePath,
		SSHClient:  sshClient,
		SFTPClient: sftpClient,
	}, nil
}

func (repo *SSHRepository) Upload(localPath, repoPath string, m *Meter) error {
	client := repo.SFTPClient

	sourceFileStat, err := os.Stat(localPath)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", localPath)
	}

	source, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer source.Close()

	// Copy from source to tmp
	tmpDir := path.Join(repo.BaseDir, "tmp")
	err = client.MkdirAll(tmpDir)
	if err != nil {
		return err
	}

	tmpPath := filepath.Join(repo.BaseDir, "tmp", fmt.Sprintf("%d", (rand.Int()%100000000)))
	tmp, err := client.Create(tmpPath)
	if err != nil {
		return err
	}
	defer client.Remove(tmpPath)
	_, err = CopyWithMeter(tmp, source, m)
	if err != nil {
		return err
	}
	err = tmp.Close()
	if err != nil {
		return err
	}

	// Move from tmp to dest
	destPath := path.Join(repo.BaseDir, repoPath)
	err = client.MkdirAll(filepath.Dir(destPath))
	if err != nil {
		return err
	}
	err = client.Remove(destPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	err = client.Rename(tmpPath, destPath)
	if err != nil {
		return err
	}

	return nil
}

func (repo *SSHRepository) Download(repoPath, localPath string, m *Meter) error {
	client := repo.SFTPClient

	srcPath := path.Join(repo.BaseDir, repoPath)
	src, err := client.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dest, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer dest.Close()
	written, err := CopyWithMeter(dest, src, m)
	if err != nil {
		return err
	}

	if written == 0 {
		err = os.Truncate(localPath, 0)
	}

	return err
}

func (repo *SSHRepository) Delete(repoPath string) error {
	filePath := path.Join(repo.BaseDir, repoPath)
	return repo.SFTPClient.Remove(filePath)
}

func (repo *SSHRepository) Stat(repoPath string) (FileInfo, error) {
	filePath := path.Join(repo.BaseDir, repoPath)
	return repo.SFTPClient.Stat(filePath)
}

func (repo *SSHRepository) List(repoPath string) ([]FileInfo, error) {
	client := repo.SFTPClient

	dir := path.Join(repo.BaseDir, repoPath)
	fs, err := client.ReadDir(dir)
	if err != nil {
		return []FileInfo{}, nil
	}
	fs2 := []FileInfo{}

	for _, info := range fs {
		info2, ok := info.(FileInfo)
		if ok {
			fs2 = append(fs2, info2)
		}
	}
	return fs2, nil
}
