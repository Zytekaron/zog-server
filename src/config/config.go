package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Cache    *Cache    `yaml:"cache"`
	Server   *Server   `yaml:"server"`
	Database *Database `yaml:"db"`
}

type Cache struct {
	Tokens int `yaml:"tokens"`
}

type Server struct {
	Domain string `yaml:"domain"`
	Addr   string `yaml:"addr"`
}

type Database struct {
	Use     string   `yaml:"use"`
	MongoDB *MongoDB `yaml:"mongodb"`
}

type MongoDB struct {
	URI             string `yaml:"uri"`
	Database        string `yaml:"database"`
	LogCollection   string `yaml:"log_collection"`
	TokenCollection string `yaml:"token_collection"`
}

func Load(dirs []string) (*Config, error) {
	var cfg *Config
	for i := 0; i < len(dirs); i++ {
		file, err := os.Open(dirs[i])
		if err != nil {
			continue
		}

		return cfg, yaml.NewDecoder(file).Decode(&cfg)
	}

	return nil, errors.New("could not find any config files")
}
