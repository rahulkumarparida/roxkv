package kvagent

import "github.com/ollama/ollama/api"


var KvInference = map[string]any{
    "num_predict": 200,   // Limit output to a maximum of 100 tokens
    "num_ctx":     1048,  // Set total context window (input + output) to 2048 tokens
    "temperature": 0.4,   // Lower temperature makes responses more focused and deterministic
    "top_p":       0.9,   // Top-p sampling boundary
}


func GetKeyTool() api.Tool{

var properties = api.NewToolPropertiesMap()

keyprop := api.ToolProperty{
	Type: api.PropertyType{"string"},
	Description: "The key is used to search through the in memopry database",
}

properties.Set("key",keyprop)

toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
		Required:   []string{"key"},
	}


var GetKeyValueTool = api.Tool{
	Type: "function",
	Function: api.ToolFunction{
		Name:        "get",
		Description: "Use this tool to retrieve, fetch, look up, or read the stored value of a specific key from the database memory.",
		Parameters:  toolParams,
	},
}
return GetKeyValueTool

}



func SetKeyTool() api.Tool{

var properties = api.NewToolPropertiesMap()

keyprop := api.ToolProperty{
	Type:        api.PropertyType{"string"},
	Description: "The unique identifier or name of the key to look up, set, or modify in the database.",
}

valprop := api.ToolProperty{
	Type:        api.PropertyType{"string"},
	Description: "The actual data payload or content to assign and store under the specified key.",
}

properties.Set("key",keyprop)
properties.Set("value",valprop)

toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
		Required:   []string{"key","value"},
	}


var SetKeyValueTool = api.Tool{
	Type: "function",
	Function: api.ToolFunction{
		Name:        "set",
		Description: "Use this tool to create, store, save, or update a key-value pair in the database memory.",
		Parameters:  toolParams,
	},
}

return SetKeyValueTool

}


func DeleteKeyTool() api.Tool {
	var properties = api.NewToolPropertiesMap()

	keyprop := api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "The unique identifier or name of the key to remove from the database.",
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
			Description: "Use this tool to delete, remove, or erase a specific key and its associated value from the database memory. Only when explicitly mentioned remove or delete or erase a key",
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
			Description: "Use this tool to list, retrieve, or fetch all existing keys, identifiers currently stored in the database memory.",
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
			Description: "Use this tool to persist, commit, write, or save the current database memory snapshot onto the disk or hard drive for data persistence.",
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
			Description: "Use this tool to read, restore, load, or recover previously persisted key-value data from the disk or hard drive back into active database memory.",
			Parameters:  toolParams,
		},
	}

	return LoadDiskTool
}
