package abstractor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/rahulkumarparida/roxkv/internal/utils"
)

var (
	cachedConfig *ProviderConfig
	configOnce   sync.Once
	configMu     sync.RWMutex
)

// DefaultConfig returns the default Ollama configuration.
func DefaultConfig() ProviderConfig {
	return ProviderConfig{
		Provider:    "ollama",
		Model:       "llama3.2:3b",
		Endpoint:    "http://localhost:11434",
		APIKey:      "",
		Temperature: 0.2,
		TopP:        0.9,
		MaxTokens:   400,
		Stream:      false,
		Timeout:     30,
	}
}

// ConfigPath returns the absolute path to model.json.
// It resolves relative to the executable's directory, walking up to find the project root.
func ConfigPath() string {
	// Try relative to current working directory first
	configwd:= utils.ConfigFolder()
	
	return filepath.Join(configwd, "model.json")
}

// LoadConfig reads model.json from disk, creating default if missing.
// The result is cached — subsequent calls return the cached value.
func LoadConfig() (*ProviderConfig, error) {
	var loadErr error
	configOnce.Do(func() {
		path := ConfigPath()

		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				// Create default config
				defCfg := DefaultConfig()
				if writeErr := writeConfigToDisk(path, &defCfg); writeErr != nil {
					loadErr = fmt.Errorf("failed to create default config: %w", writeErr)
					return
				}
				cachedConfig = &defCfg
				return
			}
			loadErr = fmt.Errorf("failed to read config: %w", err)
			return
		}
		fmt.Println("Model Data:",string(data))

		var cfg ProviderConfig
		if err := json.Unmarshal(data, &cfg); err != nil {
			loadErr = fmt.Errorf("failed to parse config: %w", err)
			return
		}
		cachedConfig = &cfg
	})

	if loadErr != nil {
		return nil, loadErr
	}
	return cachedConfig, nil
}

// GetConfig returns the cached config. Panics if LoadConfig was never called.
func GetConfig() *ProviderConfig {
	configMu.RLock()
	defer configMu.RUnlock()
	if cachedConfig == nil {
		defCfg := DefaultConfig()
		return &defCfg
	}
	copycon := *cachedConfig
	return &copycon
}

// SaveConfig writes the config to disk and updates the cache.
func SaveConfig(cfg ProviderConfig) error {
	configMu.Lock()
	defer configMu.Unlock()

	path := ConfigPath()
	if err := writeConfigToDisk(path, &cfg); err != nil {
		return err
	}
	cachedConfig = &cfg
	return nil
}

// ReloadConfig forces a re-read of the config from disk.
func ReloadConfig() (*ProviderConfig, error) {
	configMu.Lock()
	defer configMu.Unlock()

	path := ConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg ProviderConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	cachedConfig = &cfg
	return cachedConfig, nil
}

func writeConfigToDisk(path string, cfg *ProviderConfig) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	return nil
}
