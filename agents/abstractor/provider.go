package abstractor

import "context"

// Provider is the interface every LLM backend must implement.
// The Master Agent communicates exclusively through this interface.
// Providers NEVER execute tools — they only return requested tool calls.
type Provider interface {
	// Chat sends a conversation with optional tools to the LLM and returns
	// either text, tool calls, or both. The provider handles all SDK-specific
	// serialization and deserialization internally.
	Chat(ctx context.Context, messages []GenericMessage, tools []GenericToolDefinition, config ProviderConfig) (*ProviderResponse, error)

	// Name returns the provider identifier (e.g., "ollama", "openai").
	Name() string
}
