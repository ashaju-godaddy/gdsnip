package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// Config holds CLI configuration
type Config struct {
	APIURL string `koanf:"api_url"`
}

// Load loads CLI configuration from file and environment
func Load() (*Config, error) {
	k := koanf.New(".")

	// Get config file path: ~/.gdsnip/config.yml
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	// Load from config file if exists (optional)
	if _, err := os.Stat(configPath); err == nil {
		if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
			return nil, fmt.Errorf("error loading config file: %w", err)
		}
	}

	// Load from environment variables with GDSNIP_ prefix
	if err := k.Load(env.Provider("GDSNIP_", ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, "GDSNIP_"))
	}), nil); err != nil {
		return nil, fmt.Errorf("error loading environment variables: %w", err)
	}

	// Unmarshal into Config struct
	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Apply default if not set
	if cfg.APIURL == "" {
		cfg.APIURL = "http://localhost:8080/v1"
	}

	return &cfg, nil
}

// getConfigPath returns the path to the config file
func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get home directory: %w", err)
	}

	return filepath.Join(home, ".gdsnip", "config.yml"), nil
}

// GetConfigDir returns the ~/.gdsnip directory path
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get home directory: %w", err)
	}

	return filepath.Join(home, ".gdsnip"), nil
}
