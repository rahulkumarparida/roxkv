package kvagent

import "github.com/ollama/ollama/api"

var KvInference = map[string]any{
	"num_predict": 200,  // Limit output to a maximum of 100 tokens
	"num_ctx":     1048, // Set total context window (input + output) to 2048 tokens
	"temperature": 0.4,  // Lower temperature makes responses more focused and deterministic
	"top_p":       0.9,  // Top-p sampling boundary
}

func GetKeyTool() api.Tool {

	var properties = api.NewToolPropertiesMap()

	keyprop := api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Exact key name to read from the in-memory store. Use a single stored key identifier.",
	}

	properties.Set("key", keyprop)

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
		Required:   []string{"key"},
	}

	var GetKeyValueTool = api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get",
			Description: "Purpose: read one stored value by key. Inputs: key string. Output: the value currently stored for that key. Use when the request asks for an exact key lookup. Do not use for listing keys, writing data, or storage analytics.",
			Parameters:  toolParams,
		},
	}
	return GetKeyValueTool

}

func SetKeyTool() api.Tool {

	var properties = api.NewToolPropertiesMap()

	keyprop := api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Exact key name to create or overwrite in the in-memory store.",
	}

	valprop := api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "String value to store under the provided key.",
	}

	properties.Set("key", keyprop)
	properties.Set("value", valprop)

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
		Required:   []string{"key", "value"},
	}

	var SetKeyValueTool = api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "set",
			Description: "Purpose: create or overwrite one key-value pair in memory. Inputs: key string, value string. Output: boolean-style write result from the store. Use when the request explicitly asks to save, set, or update a key. Do not use for reads, deletes, or disk persistence.",
			Parameters:  toolParams,
		},
	}

	return SetKeyValueTool

}

func DeleteKeyTool() api.Tool {
	var properties = api.NewToolPropertiesMap()

	keyprop := api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Exact key name to remove from the in-memory store.",
	}

	properties.Set("key", keyprop)

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
		Required:   []string{"key"},
	}

	var DeleteValueTool = api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "del",
			Description: "Purpose: delete one key and its value from memory. Inputs: key string. Output: boolean-style deletion result. Use when the request explicitly asks to delete or remove a key. Do not use for reads, listings, or persistence.",
			Parameters:  toolParams,
		},
	}

	return DeleteValueTool
}

func KeysTool() api.Tool {
	var properties = api.NewToolPropertiesMap()

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
		Required:   []string{}, // No required fields since it lists everything
	}

	var ListKeysTool = api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "keys",
			Description: "Purpose: list every key currently present in the in-memory store. Inputs: none. Output: array of key names. Use when the request asks to enumerate keys. Do not use for reading one value, writing data, or metadata analysis.",
			Parameters:  toolParams,
		},
	}

	return ListKeysTool
}

func SaveTool() api.Tool {
	var properties = api.NewToolPropertiesMap()

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
		Required:   []string{},
	}

	var SaveDiskTool = api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "save",
			Description: "Purpose: persist the current in-memory database to disk. Inputs: none. Output: save status string from the persistence command. Use when the request asks to save current state. Do not use for reading from disk or inspecting snapshot health.",
			Parameters:  toolParams,
		},
	}

	return SaveDiskTool
}

func LoadTool() api.Tool {
	var properties = api.NewToolPropertiesMap()

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
		Required:   []string{},
	}

	var LoadDiskTool = api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "load",
			Description: "Purpose: load previously persisted key-value data from disk into memory. Inputs: none. Output: number of records loaded. Use when the request asks to restore saved data. Do not use for saving, reading one key, or listing keys already in memory.",
			Parameters:  toolParams,
		},
	}

	return LoadDiskTool
}
