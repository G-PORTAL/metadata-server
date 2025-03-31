package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

const defaultConfigPath = "/etc/metadata-server/config.yaml"

var cfg *Config

var errInvalidConfig = errors.New("invalid config")

func ReloadConfig() error {
	cfg = &Config{}
	cfg.loadDefaults()

	content, err := getConfigContent()
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(content, cfg); err != nil {
		return errInvalidConfig
	}

	return nil
}

func GetConfig() *Config {
	if cfg == nil {
		_ = ReloadConfig()
	}

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
