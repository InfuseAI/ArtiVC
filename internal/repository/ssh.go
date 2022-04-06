package repository

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/infuseai/artivc/internal/log"
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

func normalizeKeyPath(path string) string {
	if path == "~" {
		usr, _ := user.Current()
		return usr.HomeDir
	} else if strings.HasPrefix(path, "~/") {
		usr, _ := user.Current()
		return filepath.Join(usr.HomeDir, path[2:])
	} else {
		return path
	}
}

func NewSSHRepository(host, basePath string) (*SSHRepository, error) {

	explictSigners := []ssh.Signer{}
	var agentClient agent.ExtendedAgent

	authPublickey := ssh.PublicKeysCallback(func() (signers []ssh.Signer, err error) {
		if agentClient == nil {
			return explictSigners, nil
		}

		agentSigners, err := agentClient.Signers()
		if err != nil {
			return agentSigners, err
		}
		return append(agentSigners, explictSigners...), nil
	})

	authMethods := []ssh.AuthMethod{authPublickey}
	agentSock := os.Getenv("SSH_AUTH_SOCK")
	if agentSock != "" {
		agentConn, err := net.Dial("unix", agentSock)
		if err != nil {
			return nil, err
		}
		agentClient = agent.NewClient(agentConn)
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
	identifierFiles, _ := cfg.GetAll(host, "IdentityFile")
	passphrase := os.Getenv("SSH_KEY_PASSPHRASE")
	fmt.Println(user)

	for _, identityFile := range identifierFiles {
		key, err := ioutil.ReadFile(normalizeKeyPath(identityFile))
		if err != nil {
			log.Debugf("cannot parse key %s: %s", identityFile, err.Error())
			continue
		}

		var signer ssh.Signer
		if passphrase == "" {
			signer, err = ssh.ParsePrivateKey(key)
			if err != nil {
				log.Debugf("cannot parse key %s: %s", identityFile, err.Error())
				continue
			}
		} else {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(key, []byte(passphrase))
			if err != nil {
				log.Debugf("cannot parse key %s: %s", identityFile, err.Error())
				continue
			}
		}

		explictSigners = append(explictSigners, signer)
	}

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

	sftpClient, err := sftp.NewClient(sshClient, sftp.UseConcurrentReads(true), sftp.UseConcurrentWrites(true))
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

	_, err = tmp.ReadFrom(&sshFileWrapper{file: source, meter: m})
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

	written, err := src.WriteTo(&sshFileWrapper{file: dest, meter: m})
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

type sshFileWrapper struct {
	file  *os.File
	meter *Meter
}

func (r *sshFileWrapper) Read(p []byte) (n int, err error) {
	n, err = r.file.Read(p)
	if err == nil && r.meter != nil {
		r.meter.AddBytes(n)
	}
	return
}

func (r *sshFileWrapper) Write(p []byte) (n int, err error) {
	n, err = r.file.Write(p)
	if err == nil && r.meter != nil {
		r.meter.AddBytes(n)
	}
	return
}

func (r *sshFileWrapper) Stat() (os.FileInfo, error) {
	return r.file.Stat()
}

func (r *sshFileWrapper) Close() error {
	return r.file.Close()
}
