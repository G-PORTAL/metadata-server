package config

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"os"
)

const defaultConfigPath = "/etc/metadata-server/config.yaml"

var cfg *Config

var errInvalidConfig = fmt.Errorf("invalid config")

func ReloadConfig() error {
	cfg = &Config{}
	cfg.loadDefaults()

	content, err := getConfigContent()
	if err != nil {
		return err
	}

	if err = yaml.UnmarshalStrict(content, cfg); err != nil {
		return errInvalidConfig
	}

	return nil
}

func GetConfig() *Config {
	return cfg
}

func getConfigContent() ([]byte, error) {
	data, err := os.ReadFile(getConfigFilePath())
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	return data, nil
}

func getConfigFilePath() string {
	if path := os.Getenv("METADATA_SERVER_CONFIG"); path != "" {
		return path
	}

	return defaultConfigPath
}
