package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// Config holds all API configuration
type Config struct {
	Port        string        `koanf:"port"`
	DatabaseURL string        `koanf:"database_url"`
	JWTSecret   string        `koanf:"jwt_secret"`
	JWTExpiry   time.Duration `koanf:"jwt_expiry"`
	MaxBodySize string        `koanf:"max_body_size"`
	RateLimit   int           `koanf:"rate_limit"`
}

// Load loads configuration from file and environment variables
func Load(configFile string) (*Config, error) {
	k := koanf.New(".")

	// Load .env file if it exists (silently ignore if not present)
	_ = godotenv.Load()

	// Load from config file if exists
	if configFile != "" {
		if err := k.Load(file.Provider(configFile), yaml.Parser()); err != nil {
			// Config file is optional, only error if file exists but can't be parsed
			return nil, fmt.Errorf("error loading config file: %w", err)
		}
	}

	// Load from environment variables with GDSNIP_ prefix
	// Environment variables override config file
	if err := k.Load(env.Provider("GDSNIP_", ".", func(s string) string {
		// GDSNIP_DATABASE_URL -> database_url
		return strings.ToLower(strings.TrimPrefix(s, "GDSNIP_"))
	}), nil); err != nil {
		return nil, fmt.Errorf("error loading environment variables: %w", err)
	}

	// Unmarshal into Config struct
	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Apply defaults if not set
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.MaxBodySize == "" {
		cfg.MaxBodySize = "1MB"
	}
	if cfg.RateLimit == 0 {
		cfg.RateLimit = 100
	}

	// Parse JWT expiry string to duration
	expiryStr := k.String("jwt_expiry")
	if expiryStr != "" {
		duration, err := time.ParseDuration(expiryStr)
		if err != nil {
			return nil, fmt.Errorf("invalid jwt_expiry format: %w", err)
		}
		cfg.JWTExpiry = duration
	} else {
		// Default to 7 days
		cfg.JWTExpiry = 168 * time.Hour
	}

	// Validate required fields
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate checks that required configuration is present
func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("database_url is required (set GDSNIP_DATABASE_URL)")
	}

	if c.JWTSecret == "" {
		return fmt.Errorf("jwt_secret is required (set GDSNIP_JWT_SECRET)")
	}

	if len(c.JWTSecret) < 32 {
		return fmt.Errorf("jwt_secret must be at least 32 characters")
	}

	if c.Port == "" {
		return fmt.Errorf("port is required")
	}

	return nil
}
