package abstractor

import "fmt"

// NewProvider creates a Provider implementation based on the config's Provider field.
func NewProvider(config ProviderConfig) (Provider, error) {
	switch config.Provider {
	case "ollama":
		return NewOllamaProvider(config)
	case "openai":
		return NewOpenAIProvider(config)
	case "claude":
		return NewClaudeProvider(config)
	case "gemini":
		return NewGeminiProvider(config)
	case "groq":
		return NewGroqProvider(config)
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownProvider, config.Provider)
	}
}
