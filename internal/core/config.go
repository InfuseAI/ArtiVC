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

func InitRepo(baseDir, repo string) error {
	config := map[string]interface{}{
		"repo": map[string]interface{}{
			"url": repo,
		},
	}

	airDir := path.Join(baseDir, ".art")
	configPath := path.Join(airDir, "config")
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
	config  map[string]interface{}
	path    string
	baseDir string
}

func LoadConfig() (*ArtConfig, error) {

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

	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	for {
		config, err := load(dir)
		if errors.As(err, &toml.ParseError{}) {
			fmt.Fprintf(os.Stderr, "cannot load the workspace config\n")
			return nil, err
		}

		if err == nil {
			return &ArtConfig{config: config, baseDir: dir, path: path.Join(dir, ".art/config")}, nil
		}

		newDir := filepath.Dir(dir)
		if dir == newDir {
			break
		}
		dir = newDir
	}

	return nil, fmt.Errorf("cannot find art repository")
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

func (config *ArtConfig) BaseDir() string {
	return config.baseDir
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
	f, err := os.OpenFile(config.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
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
