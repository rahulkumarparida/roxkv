package agents

import (
	"encoding/json"
	"fmt"
)

type ToolResult struct {
	Tool   string `json:"tool"`
	Output any    `json:"output"`
}

func StructuredSuccess(results []ToolResult) string {
	return marshal(map[string]any{
		"status":  "success",
		"count":   len(results),
		"results": results,
	})
}

func StructuredError(code, message string) string {
	return marshal(map[string]any{
		"status":  "error",
		"error":   code,
		"message": message,
	})
}

func ParseToolOutput(raw string) any {
	var decoded any
	if err := json.Unmarshal([]byte(raw), &decoded); err == nil {
		return decoded
	}

	return raw
}

func marshal(value any) string {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", value)
	}

	return string(data)
}
