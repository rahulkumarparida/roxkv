package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var configMu sync.RWMutex

// ProviderConfig represents the generic configuration structure for any AI provider.
// Fields that are not applicable to a specific provider will remain empty or hold default values.
type ProviderConfig struct {
	Provider     string                 `json:"provider"`
	Model        string                 `json:"model"`
	Endpoint     string                 `json:"endpoint"`
	APIKey       string                 `json:"apiKey"`
	Organization string                 `json:"organization"`
	Temperature  float64                `json:"temperature"`
	TopP         float64                `json:"topP"`
	MaxTokens    int                    `json:"maxTokens"`
	Timeout      int                    `json:"timeout"`
	Stream       bool                   `json:"stream"`
	Headers      map[string]string      `json:"headers"`
	Extra        map[string]interface{} `json:"extra"`
}

// ConfigDirectory returns the absolute path to the roxconfig directory
func ConfigDirectory() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(homePath, ".roxkv", "roxconfig"), nil
}

// EnsureConfigDirectory ensures that the .roxkv/roxconfig directory exists.
// It creates the directory if it does not exist.
func EnsureConfigDirectory() error {
	dir, err := ConfigDirectory()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", dir, err)
	}
	return nil
}

// GetConfigPath returns the full file path for a given provider's JSON configuration.
func GetConfigPath(provider string) (string, error) {
	dir, err := ConfigDirectory()
	if err != nil {
		return "", err
	}
	// Sanitize provider name to avoid path traversal
	provider = strings.ToLower(strings.TrimSpace(provider))
	if provider == "" || strings.Contains(provider, "..") || strings.Contains(provider, "/") {
		return "", fmt.Errorf("invalid provider name: %s", provider)
	}
	return filepath.Join(dir,fmt.Sprintf("%v",provider)), nil
}

// SaveProvider saves the given ProviderConfig to the filesystem.
func SaveProvider(provider string, config ProviderConfig) error {
	configMu.Lock()
	defer configMu.Unlock()

	if err := EnsureConfigDirectory(); err != nil {
		return err
	}

	path, err := GetConfigPath(provider)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal config for %s: %w", provider, err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config to %s: %w", path, err)
	}

	return nil
}

// LoadProvider loads the ProviderConfig for the given provider from the filesystem.
func LoadProvider(provider string) (*ProviderConfig, error) {
	configMu.RLock()
	defer configMu.RUnlock()

	path, err := GetConfigPath(provider)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("provider configuration not found: %s", provider)
		}
		return nil, fmt.Errorf("failed to read config file for %s: %w", provider, err)
	}

	var config ProviderConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON config for %s: %w", provider, err)
	}
	return &config, nil
}

// DeleteProvider deletes the configuration file for the given provider.
func DeleteProvider(provider string) error {
	configMu.Lock()
	defer configMu.Unlock()

	path, err := GetConfigPath(provider)
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return nil // Already deleted
		}
		return fmt.Errorf("failed to delete config file for %s: %w", provider, err)
	}

	return nil
}

// ProviderExists checks whether a configuration file exists for the given provider.
func ProviderExists(provider string) bool {
	configMu.RLock()
	defer configMu.RUnlock()

	path, err := GetConfigPath(provider)
	if err != nil {
		return false
	}

	_, err = os.Stat(path)
	return err == nil
}

// ListConfiguredProviders returns a list of all provider names that have a configuration file.
func ListConfiguredProviders() ([]string, error) {
	configMu.RLock()
	defer configMu.RUnlock()

	dir, err := ConfigDirectory()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read config directory: %w", err)
	}

	var providers []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			name := strings.TrimSuffix(entry.Name(), ".json")
			providers = append(providers, name)
		}
	}

	return providers, nil
}
