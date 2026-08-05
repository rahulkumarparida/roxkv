package registry

import (
	"sync"

	"github.com/rahulkumarparida/roxkv/agents/abstractor"
)

// ToolMetadata describes a single tool for metadata-driven routing.
// The router scores each tool against the user query using the fields
// below, then only the highest-scoring subset is exposed to the LLM.
type ToolMetadata struct {
	// Name is the unique identifier matching api.ToolFunction.Name.
	Name string `json:"name"`

	// Tool is the actual abstractor.GenericToolDefinition value passed to the LLM when selected.
	Tool abstractor.GenericToolDefinition `json:"tool"`

	// Category is a logical grouping such as "KV", "Monitoring",
	// "PubSub", "Storage", or "Persistence".
	Category string `json:"category"`

	// Description is a concise explanation used only for routing
	// and maintainability. It is NOT sent to the LLM.
	Description string `json:"description"`

	// Keywords are primary routing keywords that trigger selection.
	Keywords []string `json:"keywords"`

	// Synonyms are natural-language alternatives for the tool's action.
	Synonyms []string `json:"synonyms"`

	// Examples are representative user request phrases used for
	// example-phrase matching during routing.
	Examples []string `json:"examples"`

	// Priority is a numeric routing weight. Higher values indicate
	// stronger confidence when a match occurs.
	Priority int `json:"priority"`

	// IsMutation indicates whether the tool changes application state.
	// Read-only tools should be false.
	IsMutation bool `json:"isMutation"`
}

var (
	mu       sync.RWMutex
	registry []ToolMetadata
)

// Register adds one tool's metadata to the global registry.
func Register(meta ToolMetadata) {
	mu.Lock()
	defer mu.Unlock()
	registry = append(registry, meta)
}

// AllMetadata returns a copy of all registered tool metadata.
func AllMetadata() []ToolMetadata {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]ToolMetadata, len(registry))
	copy(out, registry)
	return out
}

// AllTools returns every registered abstractor.GenericToolDefinition value.
// This is the fallback path when routing cannot narrow the selection.
func AllTools() []abstractor.GenericToolDefinition {
	mu.RLock()
	defer mu.RUnlock()
	tools := make([]abstractor.GenericToolDefinition, 0, len(registry))
	for _, m := range registry {
		tools = append(tools, m.Tool)
	}
	return tools
}
