package abstractor

import "strings"

// buildOpenAIToolJSON converts a GenericToolDefinition to the OpenAI/Groq tool format.
// Returns a map suitable for JSON marshalling.
func buildOpenAIToolJSON(tool GenericToolDefinition) map[string]any {
	props := make(map[string]any)
	for name, prop := range tool.Parameters.Properties {
		p := map[string]any{
			"type":        prop.Type,
			"description": prop.Description,
		}
		if len(prop.Enum) > 0 {
			p["enum"] = prop.Enum
		}
		props[name] = p
	}

	params := map[string]any{
		"type":       tool.Parameters.Type,
		"properties": props,
	}
	if len(tool.Parameters.Required) > 0 {
		params["required"] = tool.Parameters.Required
	}

	return map[string]any{
		"type": "function",
		"function": map[string]any{
			"name":        tool.Name,
			"description": tool.Description,
			"parameters":  params,
		},
	}
}

// buildClaudeToolJSON converts a GenericToolDefinition to the Claude tool format.
func buildClaudeToolJSON(tool GenericToolDefinition) map[string]any {
	// Claude uses "input_schema" instead of "parameters"
	props := make(map[string]any)
	for name, prop := range tool.Parameters.Properties {
		p := map[string]any{
			"type":        prop.Type,
			"description": prop.Description,
		}
		if len(prop.Enum) > 0 {
			p["enum"] = prop.Enum
		}
		props[name] = p
	}

	schema := map[string]any{
		"type":       tool.Parameters.Type,
		"properties": props,
	}
	if len(tool.Parameters.Required) > 0 {
		schema["required"] = tool.Parameters.Required
	}

	return map[string]any{
		"name":         tool.Name,
		"description":  tool.Description,
		"input_schema": schema,
	}
}

// buildGeminiToolJSON converts a slice of GenericToolDefinition to Gemini's tool format.
func buildGeminiToolJSON(tools []GenericToolDefinition) map[string]any {
	declarations := make([]map[string]any, 0, len(tools))
	for _, tool := range tools {
		props := make(map[string]any)
		for name, prop := range tool.Parameters.Properties {
			p := map[string]any{
				"type":        strings.ToUpper(prop.Type),
				"description": prop.Description,
			}
			if len(prop.Enum) > 0 {
				p["enum"] = prop.Enum
			}
			props[name] = p
		}

		params := map[string]any{
			"type":       "OBJECT",
			"properties": props,
		}
		if len(tool.Parameters.Required) > 0 {
			params["required"] = tool.Parameters.Required
		}

		declarations = append(declarations, map[string]any{
			"name":        tool.Name,
			"description": tool.Description,
			"parameters":  params,
		})
	}

	return map[string]any{
		"functionDeclarations": declarations,
	}
}
