package repository

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	osuser "os/user"
	"path"
	"path/filepath"
	"strconv"
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
	BaseDir    string
	SSHClient  *ssh.Client
	SFTPClient *sftp.Client
}

func normalizeKeyPath(path string) string {
	if path == "~" {
		usr, _ := osuser.Current()
		return usr.HomeDir
	} else if strings.HasPrefix(path, "~/") {
		usr, _ := osuser.Current()
		return filepath.Join(usr.HomeDir, path[2:])
	} else {
		return path
	}
}

func NewSSHRepository(hostname, basePath string) (*SSHRepository, error) {
	currentUser, err := osuser.Current()
	if err != nil {
		return nil, err
	}

	user := currentUser.Username
	port := 22

	explictSigners := []ssh.Signer{}

	// ssh config
	f, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "config"))
	if err == nil {
		cfg, err := ssh_config.Decode(f)
		if err != nil {
			return nil, err
		}

		alias := hostname
		if value, err := cfg.Get(alias, "Hostname"); err == nil {
			hostname = value
		}

		if value, err := cfg.Get(alias, "Port"); err == nil {
			port, err = strconv.Atoi(value)
			if err != nil {
				return nil, err
			}
		}

		if value, err := cfg.Get(alias, "User"); err == nil {
			user = value
		}

		if identifierFiles, err := cfg.GetAll(alias, "IdentityFile"); err == nil {
			for _, identityFile := range identifierFiles {
				signer, err := sshLoadIdentifyFile(identityFile)
				if err != nil {
					log.Debugf("cannot parse key %s: %s", identityFile, err.Error())
					continue
				}

				log.Debugln("Add identify file explict: " + identityFile)
				explictSigners = append(explictSigners, signer)
			}
		}
	}

	// SSH Agent
	var agentClient agent.ExtendedAgent
	if agentSock := os.Getenv("SSH_AUTH_SOCK"); agentSock != "" {
		agentConn, err := net.Dial("unix", agentSock)
		if err != nil {
			return nil, err
		}
		agentClient = agent.NewClient(agentConn)
	}

	// Auth method: Password
	authMethods := []ssh.AuthMethod{}
	if password := os.Getenv("SSH_PASSWORD"); password != "" {
		authMethods = append(authMethods, ssh.Password(password))
	}

	// Auth method: Public Keys
	if identityFile := os.Getenv("SSH_IDENTIFY_FILE"); identityFile != "" {
		signer, err := sshLoadIdentifyFile(identityFile)
		if err != nil {
			return nil, err
		}

		log.Debugln("Add identify file explict: " + identityFile)
		explictSigners = append(explictSigners, signer)
	}

	if agentClient != nil || len(explictSigners) > 0 {
		authPublickey := ssh.PublicKeysCallback(func() (signers []ssh.Signer, err error) {
			if agentClient == nil {
				return explictSigners, nil
			}

			agentSigners, err := agentClient.Signers()
			if err != nil {
				log.Debug("request signers from agent failed: " + err.Error())
				return explictSigners, nil
			}

			return append(agentSigners, explictSigners...), nil
		})

		authMethods = append(authMethods, authPublickey)
	}

	// setup the ssh client and sftp client
	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port), &ssh.ClientConfig{
		User:            user,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return nil, err
	}

	sftpClient, err := sftp.NewClient(sshClient, sftp.UseConcurrentReads(true), sftp.UseConcurrentWrites(true))
	if err != nil {
		return nil, err
	}

	rand.Seed(time.Now().UnixNano())
	return &SSHRepository{
		BaseDir:    basePath,
		SSHClient:  sshClient,
		SFTPClient: sftpClient,
	}, nil
}

func sshLoadIdentifyFile(identityFile string) (ssh.Signer, error) {
	key, err := ioutil.ReadFile(normalizeKeyPath(identityFile))
	if err != nil {
		return nil, err
	}

	var signer ssh.Signer
	passphrase := os.Getenv("SSH_KEY_PASSPHRASE")
	if passphrase == "" {
		signer, err = ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, err
		}
	} else {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(key, []byte(passphrase))
		if err != nil {
			return nil, err
		}
	}
	return signer, nil
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
