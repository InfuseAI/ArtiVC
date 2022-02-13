package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type ArtConfig struct {
	Repo string    `toml:"repo"`
	S3   *S3Config `toml:"s3"`
}

type S3Config struct {
	AwsAccessKeyId *string `toml:"aws_access_key_id"`
	AwsAccessKey   *string `toml:"aws_secret_access_key"`
	Region         *string `toml:"region"`
}

func s(str string) *string {
	return &str
}

func InitRepo(baseDir string) error {
	config := ArtConfig{}

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

func LoadConfig() (ArtConfig, error) {

	load := func(dir string) (ArtConfig, error) {
		var config ArtConfig
		configPath := path.Join(dir, ".art/config")

		data, err := ioutil.ReadFile(configPath)
		if err != nil {
			return config, err
		}

		err = toml.Unmarshal(data, &config)
		if err != nil {
			return config, err
		}

		return config, nil
	}

	dir, err := os.Getwd()
	if err != nil {
		return ArtConfig{}, err
	}

	for {
		config, err := load(dir)
		if err == nil {
			return config, nil
		}
		newDir := filepath.Dir(dir)
		if dir == newDir {
			break
		}
		dir = newDir
	}

	return ArtConfig{}, fmt.Errorf("cannot find art repository")
}
