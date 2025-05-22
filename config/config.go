package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds all application configuration
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

// Load loads configuration from a YAML file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}
