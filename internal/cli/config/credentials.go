package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Credentials represents stored authentication credentials
type Credentials struct {
	Token    string    `json:"token"`
	UserID   string    `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	SavedAt  time.Time `json:"saved_at"`
}

// SaveCredentials saves credentials to ~/.gdsnip/credentials.json
func SaveCredentials(creds *Credentials) error {
	// Get config directory
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist (with 0700 permissions)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal credentials to JSON
	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	// Write to file with restricted permissions (0600 - owner read/write only)
	credPath := filepath.Join(configDir, "credentials.json")
	if err := os.WriteFile(credPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write credentials file: %w", err)
	}

	return nil
}

// LoadCredentials loads credentials from ~/.gdsnip/credentials.json
func LoadCredentials() (*Credentials, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}

	credPath := filepath.Join(configDir, "credentials.json")

	// Check if file exists
	if _, err := os.Stat(credPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("not logged in: credentials file not found")
	}

	// Read file
	data, err := os.ReadFile(credPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials file: %w", err)
	}

	// Unmarshal JSON
	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("failed to parse credentials file: %w", err)
	}

	return &creds, nil
}

// ClearCredentials removes the credentials file
func ClearCredentials() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	credPath := filepath.Join(configDir, "credentials.json")

	// Remove file if it exists
	if err := os.Remove(credPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove credentials file: %w", err)
	}

	return nil
}

// GetToken returns the stored JWT token
func GetToken() (string, error) {
	creds, err := LoadCredentials()
	if err != nil {
		return "", err
	}

	return creds.Token, nil
}

// IsLoggedIn checks if user is logged in
func IsLoggedIn() bool {
	_, err := LoadCredentials()
	return err == nil
}
