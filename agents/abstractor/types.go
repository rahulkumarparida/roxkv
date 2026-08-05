package abstractor

// GenericMessage represents a single message in a conversation,
// independent of any LLM provider's SDK types.
type GenericMessage struct {
	Role      string            `json:"role"`      // "system", "user", "assistant", "tool"
	Content   string            `json:"content"`
	ToolCalls []GenericToolCall `json:"tool_calls,omitempty"` // only set when Role == "assistant"
}

// GenericToolCall represents a tool invocation request returned by an LLM.
type GenericToolCall struct {
	ID        string         `json:"id,omitempty"` // provider-assigned call ID (empty for Ollama)
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

// GenericToolDefinition describes a tool's schema (what the tool IS),
// NOT an invocation. This replaces api.Tool from the Ollama SDK.
type GenericToolDefinition struct {
	Type        string                `json:"type"`        // "function"
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Parameters  GenericToolParameters `json:"parameters"`
}

// GenericToolParameters describes the input schema for a tool.
type GenericToolParameters struct {
	Type       string                         `json:"type"` // always "object"
	Properties map[string]GenericToolProperty `json:"properties"`
	Required   []string                       `json:"required,omitempty"`
}

// GenericToolProperty describes a single property in a tool's parameter schema.
type GenericToolProperty struct {
	Type        string   `json:"type"`        // "string", "integer", "number", "boolean", "array"
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}

// ProviderResponse encapsulates everything a provider returns after a Chat call.
type ProviderResponse struct {
	Text      string
	ToolCalls []GenericToolCall
	Error     error
}

// ProviderConfig holds LLM provider settings loaded from model.json.
type ProviderConfig struct {
	Provider    string  `json:"provider"`
	Model       string  `json:"model"`
	Endpoint    string  `json:"endpoint"`
	APIKey      string  `json:"apikey"`
	Temperature float64 `json:"temperature"`
	TopP        float64 `json:"top_p"`
	MaxTokens   int     `json:"max_tokens"`
	Stream      bool    `json:"stream"`
	Timeout     int     `json:"timeout"` // seconds
}
