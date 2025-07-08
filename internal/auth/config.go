// File: internal/auth/config.go
// Purpose: Manage user config file (username, hashed password, API token)

package auth

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type UserConfig struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	APIToken     string `json:"api_token"`
}

var configPath = filepath.Join(userConfigDir(), "tor-admin", "config.json")

func userConfigDir() string {
	if dir, err := os.UserConfigDir(); err == nil {
		return dir
	}
	return "."
}

// LoadOrInitConfig loads the config or returns nil if setup is needed.
func LoadOrInitConfig() (*UserConfig, error) {
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		return nil, nil // Trigger setup flow
	}
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cfg UserConfig
	err = json.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// SaveConfig writes the config to disk.
func SaveConfig(cfg *UserConfig) error {
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	f, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(cfg)
}
