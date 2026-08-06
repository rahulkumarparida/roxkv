package abstractor

import (
	"strings"
	"sync"

	"github.com/rahulkumarparida/roxkv/internal/config"
)

// SupportedProvidersAndModels defines the central model registry for all AI providers.
var SupportedProvidersAndModels = map[string][]string{
	"ollama": {"llama3.2:3b", "llama3.2", "llama3.3", "gemma3", "deepseek"},
	"openai": {"gpt-4o", "gpt-4o-mini", "o1", "o3-mini"},
	"gemini": {"gemini-3.5-flash", "gemini-3.5-pro", "gemini-2.5-flash", "gemini-2.5-pro"},
	"claude": {"claude-3-5-sonnet-20241022", "claude-3-5-haiku-20241022", "sonnet", "opus"},
	"groq":   {"llama-3.3-70b-versatile", "llama", "mixtral"},
}

var (
	queryCountMu sync.RWMutex
	queryCounts  = make(map[string]int64)
)

// GetAvailableProviders returns the list of supported providers in consistent order.
func GetAvailableProviders() []string {
	order := []string{"ollama", "openai", "gemini", "claude", "groq"}
	var providers []string
	for _, p := range order {
		if _, exists := SupportedProvidersAndModels[p]; exists {
			providers = append(providers, p)
		}
	}
	return providers
}

// GetSupportedModels returns the models supported by a specific provider.
func GetSupportedModels(provider string) ([]string, bool) {
	models, exists := SupportedProvidersAndModels[strings.ToLower(strings.TrimSpace(provider))]
	return models, exists
}

// RequiresAPIKey returns true if the provider requires an API key for requests.
func RequiresAPIKey(provider string) bool {
	return strings.ToLower(strings.TrimSpace(provider)) != "ollama"
}

// IncrementQueryCount increments the total query count for a provider.
func IncrementQueryCount(provider string) {
	p := strings.ToLower(strings.TrimSpace(provider))
	if p == "" {
		p = "ollama"
	}
	queryCountMu.Lock()
	defer queryCountMu.Unlock()
	queryCounts[p]++
}

// GetQueryCount returns the total query count for a provider.
func GetQueryCount(provider string) int64 {
	p := strings.ToLower(strings.TrimSpace(provider))
	if p == "" {
		p = "ollama"
	}
	queryCountMu.RLock()
	defer queryCountMu.RUnlock()
	return queryCounts[p]
}

// GetConnectionStatus determines the connection status for a provider and config.
func GetConnectionStatus(provider string, cfg *config.ProviderConfig) string {
	p := strings.ToLower(strings.TrimSpace(provider))
	if RequiresAPIKey(p) {
		if cfg == nil || strings.TrimSpace(cfg.APIKey) == "" {
			return "unconfigured"
		}
	}
	return "connected"
}

// MapToAbstractorConfig converts an internal/config.ProviderConfig to abstractor.ProviderConfig
func MapToAbstractorConfig(cfg config.ProviderConfig) ProviderConfig {
	return ProviderConfig{
		Provider:     cfg.Provider,
		Model:        cfg.Model,
		Endpoint:     cfg.Endpoint,
		APIKey:       cfg.APIKey,
		Organization: cfg.Organization,
		Temperature:  cfg.Temperature,
		TopP:         cfg.TopP,
		MaxTokens:    cfg.MaxTokens,
		Stream:       cfg.Stream,
		Timeout:      cfg.Timeout,
	}
}

// MapFromAbstractorConfig converts an abstractor.ProviderConfig to internal/config.ProviderConfig
func MapFromAbstractorConfig(cfg ProviderConfig) config.ProviderConfig {
	return config.ProviderConfig{
		Provider:     cfg.Provider,
		Model:        cfg.Model,
		Endpoint:     cfg.Endpoint,
		APIKey:       cfg.APIKey,
		Organization: cfg.Organization,
		Temperature:  cfg.Temperature,
		TopP:         cfg.TopP,
		MaxTokens:    cfg.MaxTokens,
		Stream:       cfg.Stream,
		Timeout:      cfg.Timeout,
	}
}
