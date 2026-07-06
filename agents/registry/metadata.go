package registry

import (
	"sync"

	"github.com/ollama/ollama/api"
)

// ToolMetadata describes a single tool for metadata-driven routing.
// The router scores each tool against the user query using the fields
// below, then only the highest-scoring subset is exposed to the LLM.
type ToolMetadata struct {
	// Name is the unique identifier matching api.ToolFunction.Name.
	Name string

	// Tool is the actual api.Tool value passed to the LLM when selected.
	Tool api.Tool

	// Category is a logical grouping such as "KV", "Monitoring",
	// "PubSub", "Storage", or "Persistence".
	Category string

	// Description is a concise explanation used only for routing
	// and maintainability. It is NOT sent to the LLM.
	Description string

	// Keywords are primary routing keywords that trigger selection.
	Keywords []string

	// Synonyms are natural-language alternatives for the tool's action.
	Synonyms []string

	// Examples are representative user request phrases used for
	// example-phrase matching during routing.
	Examples []string

	// Priority is a numeric routing weight. Higher values indicate
	// stronger confidence when a match occurs.
	Priority int

	// IsMutation indicates whether the tool changes application state.
	// Read-only tools should be false.
	IsMutation bool
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

// AllTools returns every registered api.Tool value.
// This is the fallback path when routing cannot narrow the selection.
func AllTools() []api.Tool {
	mu.RLock()
	defer mu.RUnlock()
	tools := make([]api.Tool, 0, len(registry))
	for _, m := range registry {
		tools = append(tools, m.Tool)
	}
	return tools
}
