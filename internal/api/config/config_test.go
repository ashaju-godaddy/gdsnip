package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Defaults(t *testing.T) {
	// Set required env vars
	os.Setenv("GDSNIP_DATABASE_URL", "postgres://test")
	os.Setenv("GDSNIP_JWT_SECRET", "this-is-a-very-long-jwt-secret-key-for-testing-purposes-only")
	defer func() {
		os.Unsetenv("GDSNIP_DATABASE_URL")
		os.Unsetenv("GDSNIP_JWT_SECRET")
	}()

	cfg, err := Load("")
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Check defaults
	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, 168*time.Hour, cfg.JWTExpiry)
	assert.Equal(t, "1MB", cfg.MaxBodySize)
	assert.Equal(t, 100, cfg.RateLimit)
}

func TestLoad_EnvironmentOverrides(t *testing.T) {
	os.Setenv("GDSNIP_PORT", "9090")
	os.Setenv("GDSNIP_DATABASE_URL", "postgres://custom")
	os.Setenv("GDSNIP_JWT_SECRET", "super-long-secret-key-that-is-at-least-32-characters-long")
	os.Setenv("GDSNIP_JWT_EXPIRY", "24h")
	os.Setenv("GDSNIP_RATE_LIMIT", "200")
	defer func() {
		os.Unsetenv("GDSNIP_PORT")
		os.Unsetenv("GDSNIP_DATABASE_URL")
		os.Unsetenv("GDSNIP_JWT_SECRET")
		os.Unsetenv("GDSNIP_JWT_EXPIRY")
		os.Unsetenv("GDSNIP_RATE_LIMIT")
	}()

	cfg, err := Load("")
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "9090", cfg.Port)
	assert.Equal(t, "postgres://custom", cfg.DatabaseURL)
	assert.Equal(t, 24*time.Hour, cfg.JWTExpiry)
	assert.Equal(t, 200, cfg.RateLimit)
}

func TestValidate_MissingDatabaseURL(t *testing.T) {
	cfg := &Config{
		JWTSecret: "super-long-secret-key-that-is-at-least-32-characters-long",
	}
	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database_url is required")
}

func TestValidate_MissingJWTSecret(t *testing.T) {
	cfg := &Config{
		DatabaseURL: "postgres://test",
	}
	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "jwt_secret is required")
}

func TestValidate_ShortJWTSecret(t *testing.T) {
	cfg := &Config{
		DatabaseURL: "postgres://test",
		JWTSecret:   "short",
	}
	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "jwt_secret must be at least 32 characters")
}

func TestValidate_Valid(t *testing.T) {
	cfg := &Config{
		Port:        "8080",
		DatabaseURL: "postgres://test",
		JWTSecret:   "super-long-secret-key-that-is-at-least-32-characters-long",
	}
	err := cfg.Validate()
	assert.NoError(t, err)
}
