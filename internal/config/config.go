package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`

	ProductService struct {
		URL   string `yaml:"url"`
		Token string `yaml:"token"`
	} `yaml:"product_service"`

	HTTPClient struct {
		Timeout    int `yaml:"timeout"`
		MaxRetries int `yaml:"max_retries"`
		Backoff    int `yaml:"backoff"`
	} `yaml:"http_client"`

	LOMS struct {
		Address string `yaml:"address"`
	} `yaml:"loms"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
