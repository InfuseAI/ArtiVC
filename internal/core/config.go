package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

func s(str string) *string {
	return &str
}

func InitWorkspace(baseDir, repo string) error {
	config := map[string]interface{}{
		"repo": map[string]interface{}{
			"url": repo,
		},
	}

	configPath := path.Join(baseDir, ".art/config")
	err := mkdirsForFile(configPath)
	if err != nil {
		return err
	}

	f, err := os.Create(configPath)
	if err != nil {
		return err
	}

	if err := toml.NewEncoder(f).Encode(config); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err

	}
	return nil
}

type ArtConfig struct {
	config      map[string]interface{}
	MetadataDir string
	BaseDir     string
}

func NewConfig(baseDir, metadataDir, repoUrl string) ArtConfig {
	config := ArtConfig{
		BaseDir:     baseDir,
		MetadataDir: metadataDir,
	}
	config.config = make(map[string]interface{})
	config.SetRepoUrl(repoUrl)
	return config
}

func LoadConfig(dir string) (ArtConfig, error) {

	load := func(dir string) (map[string]interface{}, error) {
		var config = make(map[string]interface{})
		configPath := path.Join(dir, ".art/config")

		data, err := ioutil.ReadFile(configPath)
		if err != nil {
			return nil, err
		}

		err = toml.Unmarshal(data, &config)
		if err != nil {
			return nil, err
		}

		return config, nil
	}

	if dir == "" {
		var err2 error
		dir, err2 = os.Getwd()
		if err2 != nil {
			return ArtConfig{}, err2
		}
	}

	for {
		config, err := load(dir)
		var e *toml.ParseError
		if errors.As(err, &e) {
			fmt.Fprintf(os.Stderr, "cannot load the workspace config\n")
			return ArtConfig{}, err
		}

		if err == nil {
			return ArtConfig{config: config, BaseDir: dir, MetadataDir: path.Join(dir, ".art")}, nil
		}

		newDir := filepath.Dir(dir)
		if dir == newDir {
			break
		}
		dir = newDir
	}

	err2 := &WorkspaceNotFoundError{}

	return ArtConfig{}, err2
}

func (config *ArtConfig) Set(path string, value interface{}) {
	var obj map[string]interface{} = config.config

	parts := strings.Split(path, ".")
	for i, p := range parts {
		if i == len(parts)-1 {
			obj[p] = value
		} else {
			if v, ok := obj[p].(map[string]interface{}); ok {
				obj = v
			} else {
				child := make(map[string]interface{})
				obj[p] = child
				obj = child
			}
		}
	}
}

func (config *ArtConfig) Get(path string) interface{} {
	var obj interface{} = config.config
	var val interface{} = nil

	parts := strings.Split(path, ".")
	for _, p := range parts {
		if v, ok := obj.(map[string]interface{}); ok {
			obj = v[p]
			val = obj
		} else {
			return nil
		}
	}

	return val
}

func (config *ArtConfig) GetString(path string) string {
	var value string

	if config.Get(path) != nil {
		value = config.Get(path).(string)
	}

	return value
}

func (config *ArtConfig) RepoUrl() string {
	return config.GetString("repo.url")
}

func (config *ArtConfig) SetRepoUrl(repoUrl string) {
	config.Set("repo.url", repoUrl)
}

func (config *ArtConfig) Print() {
	var printChild func(string, interface{})

	printChild = func(path string, obj interface{}) {
		if v, ok := obj.(map[string]interface{}); ok {
			for key, value := range v {
				if path == "" {
					printChild(key, value)
				} else {
					printChild(path+"."+key, value)
				}
			}
		} else {
			fmt.Printf("%s=%v\n", path, obj)
		}
	}

	printChild("", config.config)
}

func (config *ArtConfig) Save() error {
	configPath := path.Join(config.MetadataDir, "config")
	f, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	if err := toml.NewEncoder(f).Encode(config.config); err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err

	}

	return nil
}
