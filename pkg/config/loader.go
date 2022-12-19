package config

import (
	"github.com/go-yaml/yaml"
	"os"
)

const defaultConfigPath = "/etc/metadata-server/config.yaml"

var cfg *config

func ReloadConfig() error {
	cfg = &config{}
	cfg.loadDefaults()

	content, err := getConfigContent()
	if err != nil {
		return err
	}

	if err = yaml.UnmarshalStrict(content, cfg); err != nil {
		return err
	}

	return nil
}

func GetConfig() *config {
	return cfg
}

func getConfigContent() ([]byte, error) {
	return os.ReadFile(getConfigFilePath())
}

func getConfigFilePath() string {
	if path := os.Getenv("METADATA_SERVER_CONFIG"); path != "" {
		return path
	}

	return defaultConfigPath
}
