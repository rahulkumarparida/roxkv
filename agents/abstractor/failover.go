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
	"github.com/rahulkumarparida/roxkv/internal/logger"
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

func applyOptions(cfg *ProviderConfig, options map[string]any) {
	if options == nil {
		return
	}
	if temp, ok := options["temperature"]; ok {
		switch v := temp.(type) {
		case float64:
			cfg.Temperature = v
		case float32:
			cfg.Temperature = float64(v)
		case int:
			cfg.Temperature = float64(v)
		}
	}
	if topP, ok := options["top_p"]; ok {
		switch v := topP.(type) {
		case float64:
			cfg.TopP = v
		case float32:
			cfg.TopP = float64(v)
		case int:
			cfg.TopP = float64(v)
		}
	}
	if maxTokens, ok := options["max_tokens"]; ok {
		switch v := maxTokens.(type) {
		case int:
			cfg.MaxTokens = v
		case float64:
			cfg.MaxTokens = int(v)
		}
	} else if numPredict, ok := options["num_predict"]; ok {
		switch v := numPredict.(type) {
		case int:
			cfg.MaxTokens = v
		case float64:
			cfg.MaxTokens = int(v)
		}
	}
}

// ExecuteWithFailover is the central entry point for LLM chat invocations with automatic retry & failover.
func ExecuteWithFailover(ctx context.Context, initialProvider Provider, model string, options map[string]any, messages []GenericMessage, tools []GenericToolDefinition) (*ProviderResponse, string, error) {
	var initialActive string
	if initialProvider != nil && strings.TrimSpace(initialProvider.Name()) != "" {
		initialActive = strings.ToLower(strings.TrimSpace(initialProvider.Name()))
	} else {
		activeCfg := GetConfig()
		if activeCfg != nil {
			initialActive = strings.ToLower(strings.TrimSpace(activeCfg.Provider))
		}
	}

	sequence := GetFailoverSequence(initialActive)
	if len(sequence) == 0 {
		msg := "All configured LLM providers are currently unavailable. Please wait until the services are available and try again later."
		logger.ErrorLog("[LLM Manager] " + msg)
		return nil, "", &FailoverError{
			Message:            msg,
			ProvidersAttempted: []string{},
		}
	}

	var attemptedProviders []string

	for _, providerName := range sequence {
		attemptedProviders = append(attemptedProviders, providerName)
		logger.InfoLog("[LLM Manager] Attempting provider: " + providerName)

		var pConfig ProviderConfig
		activeCfg := GetConfig()
		if loaded, err := config.LoadProvider(providerName); err == nil && loaded != nil {
			pConfig = MapToAbstractorConfig(*loaded)
			if activeCfg != nil && strings.EqualFold(activeCfg.Provider, providerName) {
				if pConfig.APIKey == "" {
					pConfig.APIKey = activeCfg.APIKey
				}
				if pConfig.Endpoint == "" {
					pConfig.Endpoint = activeCfg.Endpoint
				}
				if pConfig.Model == "" {
					pConfig.Model = activeCfg.Model
				}
			}
		} else {
			if activeCfg != nil {
				pConfig = *activeCfg
			}
			pConfig.Provider = providerName
			models, ok := GetSupportedModels(providerName)
			if ok && len(models) > 0 && pConfig.Model == "" {
				pConfig.Model = models[0]
			}
			if providerName == "ollama" && pConfig.Endpoint == "" {
				pConfig.Endpoint = "http://localhost:11434"
			}
		}

		if strings.EqualFold(providerName, initialActive) && strings.TrimSpace(model) != "" {
			pConfig.Model = strings.TrimSpace(model)
		}

		applyOptions(&pConfig, options)

		var providerImpl Provider
		var err error
		if initialProvider != nil && strings.EqualFold(initialProvider.Name(), providerName) {
			providerImpl = initialProvider
		} else {
			providerImpl, err = NewProvider(pConfig)
			if err != nil {
				logger.ErrorLog("[LLM Manager] Failed to instantiate provider " + providerName + ": " + err.Error())
				continue
			}
		}

		var lastErr error
		var resp *ProviderResponse

		// Retry Loop: up to 3 retries (total 4 attempts)
		for attempt := 1; attempt <= 4; attempt++ {
			if ctx.Err() != nil {
				return nil, providerName, ctx.Err()
			}

			if attempt > 1 {
				logger.InfoLog("[LLM Manager] Attempt " + fmt.Sprintf("%d", attempt) + " for provider " + providerName)
			}

			resp, lastErr = providerImpl.Chat(ctx, messages, tools, pConfig)
			if lastErr == nil && resp != nil {
				logger.SucessLog("[LLM Manager] Provider " + providerName + " responded successfully")
				IncrementQueryCount(providerName)
				return resp, providerName, nil
			}

			logger.ErrorLog("[LLM Manager] Error from provider " + providerName + ": " + fmt.Sprintf("%v", lastErr))

			if !IsRetryableError(lastErr) {
				logger.InfoLog("[LLM Manager] Non-retryable error for " + providerName + ". Switching to next provider.")
				break
			}

			if attempt == 1 {
				logger.InfoLog("[LLM Manager] Retrying provider " + providerName + " in 1s")
				time.Sleep(1 * time.Second)
			} else if attempt == 2 {
				logger.InfoLog("[LLM Manager] Retrying provider " + providerName + " in 2s")
				time.Sleep(2 * time.Second)
			} else if attempt == 3 {
				logger.ErrorLog("[LLM Manager] Provider " + providerName + " unavailable after 3 attempts")
				break
			}
		}
	}

	allFailedMsg := "All configured LLM providers are currently unavailable. Please wait until the services are available and try again later."
	logger.ErrorLog("[LLM Manager] All providers failed. Attempted: " + fmt.Sprintf("%v", attemptedProviders))

	return nil, "", &FailoverError{
		Message:            allFailedMsg,
		ProvidersAttempted: attemptedProviders,
	}
}
