package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger LoggerConf `yaml:"logger"`
	DB     DBConf     `yaml:"db"`
	HTTP   HTTPConf   `yaml:"http"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
}

type DBConf struct {
	Type string   `yaml:"type"`
	PSQL PSQLConf `yaml:"psql"`
}

type PSQLConf struct {
	URL string `yaml:"url"`
}

type HTTPConf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func NewConfig() *Config {
	return &Config{
		Logger: LoggerConf{
			Level: "info",
		},
	}
}

func ParseConfig(filePath string) (*Config, error) {
	c := NewConfig()

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0o644)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, c); err != nil {
		return nil, err
	}

	return c, nil
}
