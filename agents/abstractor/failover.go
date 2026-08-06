package abstractor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/rahulkumarparida/roxkv/internal/config"
)

type FailoverResult struct {
	Success            bool     `json:"success"`
	Error              string   `json:"error"`
	ProvidersAttempted []string `json:"providers_attempted"`
}

type FailoverError struct {
	Message            string
	ProvidersAttempted []string
}

func (e *FailoverError) Error() string {
	res := FailoverResult{
		Success:            false,
		Error:              e.Message,
		ProvidersAttempted: e.ProvidersAttempted,
	}
	bytes, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return e.Message
	}
	return string(bytes)
}

// IsProviderConfigured checks if a provider is fully configured with credentials and model.
func IsProviderConfigured(provider string) bool {
	p := strings.ToLower(strings.TrimSpace(provider))
	if p == "" {
		return false
	}

	if p == "ollama" {
		return true // Ollama does not require an API key
	}

	loadedCfg, err := config.LoadProvider(p)
	if err != nil || loadedCfg == nil {
		return false
	}

	return strings.TrimSpace(loadedCfg.APIKey) != ""
}

// GetConfiguredProviders returns all providers that have valid configuration and credentials.
func GetConfiguredProviders() []string {
	available := GetAvailableProviders()
	var configured []string
	for _, p := range available {
		if IsProviderConfigured(p) {
			configured = append(configured, p)
		}
	}
	return configured
}

// GetFailoverSequence returns provider candidates starting with activeProvider.
func GetFailoverSequence(activeProvider string) []string {
	configured := GetConfiguredProviders()
	if len(configured) == 0 {
		return nil
	}

	active := strings.ToLower(strings.TrimSpace(activeProvider))
	var sequence []string

	if active != "" && IsProviderConfigured(active) {
		sequence = append(sequence, active)
	}

	for _, p := range configured {
		if p != active {
			sequence = append(sequence, p)
		}
	}

	return sequence
}

// IsRetryableError classifies errors to decide whether to retry the current provider or failover immediately.
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	var provErr *ProviderError
	if errors.As(err, &provErr) {
		code := provErr.StatusCode
		if code == http.StatusUnauthorized || code == http.StatusForbidden || code == http.StatusBadRequest || code == http.StatusNotFound {
			return false // Non-retryable HTTP errors
		}
		if code == http.StatusServiceUnavailable || code == http.StatusBadGateway || code == http.StatusGatewayTimeout || code == http.StatusTooManyRequests {
			return true // Retryable HTTP errors
		}
		if errors.Is(provErr.Err, ErrInvalidAPIKey) {
			return false
		}
		if errors.Is(provErr.Err, ErrConnectionRefused) || errors.Is(provErr.Err, ErrTimeout) || errors.Is(provErr.Err, ErrProviderUnavailable) || errors.Is(provErr.Err, ErrQuotaExceeded) {
			return true
		}
	}

	errStr := strings.ToLower(err.Error())
	nonRetryableTerms := []string{
		"invalid api key", "unauthorized", "401", "403", "forbidden", "404", "400", "bad request", "invalid_request",
	}
	for _, term := range nonRetryableTerms {
		if strings.Contains(errStr, term) {
			return false
		}
	}

	retryableTerms := []string{
		"503", "502", "504", "service unavailable", "bad gateway", "timeout", "connection refused", "connection reset", "eof", "temporary", "unavailable",
	}
	for _, term := range retryableTerms {
		if strings.Contains(errStr, term) {
			return true
		}
	}

	return true
}

// ActiveProviderChangeListener registers callbacks when active provider switches via failover.
var (
	onProviderChangeMu sync.RWMutex
	onProviderChange   func(newProvider, newModel string)
)

func SetOnProviderChangeListener(listener func(newProvider, newModel string)) {
	onProviderChangeMu.Lock()
	defer onProviderChangeMu.Unlock()
	onProviderChange = listener
}

func notifyProviderChange(newProvider, newModel string) {
	onProviderChangeMu.RLock()
	listener := onProviderChange
	onProviderChangeMu.RUnlock()

	if listener != nil {
		listener(newProvider, newModel)
	}
}

// ExecuteWithFailover is the central entry point for LLM chat invocations with automatic retry & failover.
func ExecuteWithFailover(ctx context.Context, messages []GenericMessage, tools []GenericToolDefinition) (*ProviderResponse, string, error) {
	activeCfg := GetConfig()
	initialActive := strings.ToLower(strings.TrimSpace(activeCfg.Provider))

	sequence := GetFailoverSequence(initialActive)
	if len(sequence) == 0 {
		msg := "All configured LLM providers are currently unavailable. Please wait until the services are available and try again later."
		fmt.Printf("[LLM Manager]\nError: %s\n", msg)
		return nil, "", &FailoverError{
			Message:            msg,
			ProvidersAttempted: []string{},
		}
	}

	var attemptedProviders []string

	for idx, providerName := range sequence {
		attemptedProviders = append(attemptedProviders, providerName)
		fmt.Printf("\n[LLM Manager]\nCurrent Provider: %s\nSending request...\n", providerName)

		pConfig := *activeCfg
		if loaded, err := config.LoadProvider(providerName); err == nil && loaded != nil {
			pConfig = MapToAbstractorConfig(*loaded)
		} else {
			pConfig.Provider = providerName
			models, ok := GetSupportedModels(providerName)
			if ok && len(models) > 0 {
				pConfig.Model = models[0]
			}
			if providerName == "ollama" && pConfig.Endpoint == "" {
				pConfig.Endpoint = "http://localhost:11434"
			}
		}

		providerImpl, err := NewProvider(pConfig)
		if err != nil {
			fmt.Printf("[LLM Manager]\nFailed to instantiate provider %s: %v\nNon-retryable error. Switching immediately to next provider...\n", providerName, err)
			continue
		}

		var lastErr error
		var resp *ProviderResponse

		// Retry Loop: up to 3 retries (total 4 attempts)
		for attempt := 1; attempt <= 4; attempt++ {
			if ctx.Err() != nil {
				return nil, providerName, ctx.Err()
			}

			if attempt > 1 {
				fmt.Printf("\nAttempt %d...\n", attempt)
			}

			resp, lastErr = providerImpl.Chat(ctx, messages, tools, pConfig)
			if lastErr == nil && resp != nil {
				fmt.Println("\nSuccess.")

				// If provider switched due to failover, persist new active provider
				if !strings.EqualFold(providerName, initialActive) {
					fmt.Printf("[LLM Manager]\nSwitching provider to %s...\nUpdating UI...\n", providerName)

					_ = SaveConfig(pConfig)
					_ = config.SaveProvider(providerName, MapFromAbstractorConfig(pConfig))
					fmt.Printf("[LLM Manager]\nCurrent Provider Updated -> %s\n", providerName)

					notifyProviderChange(providerName, pConfig.Model)
				}

				IncrementQueryCount(providerName)
				return resp, providerName, nil
			}

			fmt.Printf("\n%v\n", lastErr)

			if !IsRetryableError(lastErr) {
				fmt.Println("[LLM Manager]\nNon-retryable error.\nSwitching immediately to next provider...")
				break
			}

			if attempt == 1 {
				fmt.Println("\nRetrying in 1s...")
				time.Sleep(1 * time.Second)
			} else if attempt == 2 {
				fmt.Println("\nRetrying in 2s...")
				time.Sleep(2 * time.Second)
			} else if attempt == 3 {
				fmt.Println("\nProvider unavailable.")
				if idx < len(sequence)-1 {
					capitalizedNext := strings.Title(sequence[idx+1])
					fmt.Printf("\nSearching for next configured provider...\n\nFound: %s\n\nSwitching provider...\n", capitalizedNext)
				}
				break
			}
		}
	}

	allFailedMsg := "All configured LLM providers are currently unavailable. Please wait until the services are available and try again later."
	fmt.Printf("[LLM Manager]\nError: %s (Attempted: %v)\n", allFailedMsg, attemptedProviders)

	return nil, "", &FailoverError{
		Message:            allFailedMsg,
		ProvidersAttempted: attemptedProviders,
	}
}
