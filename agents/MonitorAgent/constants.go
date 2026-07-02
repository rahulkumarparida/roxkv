package monitoragent

import "github.com/ollama/ollama/api"

var MonitorInference = map[string]any{
	"num_predict": 400,  // Limit output to a maximum of 100 tokens
	"num_ctx":     2048, // Set total context window (input + output) to 2048 tokens
	"temperature": 0.4,  // Lower temperature makes responses more focused and deterministic
	"top_p":       0.9,  // Top-p sampling boundary
}

func GetComputerUsageTool() api.Tool {
	properties := api.NewToolPropertiesMap()

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_computer_usage",
			Description: "Returns overall computer information including username, operating system, CPU count, RAM usage, architecture, and number of running Go routines.",
			Parameters:  toolParams,
		},
	}
}

func GetCPUUsageTool() api.Tool {
	properties := api.NewToolPropertiesMap()

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_cpu_usage",
			Description: "Returns the current CPU utilization percentage of the system.",
			Parameters:  toolParams,
		},
	}
}

func GetRAMUsageTool() api.Tool {
	properties := api.NewToolPropertiesMap()

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_ram_usage",
			Description: "Returns the current RAM statistics including total memory, free memory, and memory usage percentage.",
			Parameters:  toolParams,
		},
	}
}

func GetDiskUsageTool() api.Tool {
	properties := api.NewToolPropertiesMap()

	pathProp := api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Filesystem path to inspect, for example '/' on Linux or 'C:\\' on Windows.",
	}

	properties.Set("path", pathProp)

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
		Required:   []string{"path"},
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_disk_usage",
			Description: "Returns disk statistics for the specified filesystem path including total space, free space, and available space.",
			Parameters:  toolParams,
		},
	}
}

func GetRuntimeStatsTool() api.Tool {
	properties := api.NewToolPropertiesMap()

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_runtime_stats",
			Description: "Returns Go runtime statistics including Go version, operating system, architecture, goroutine count, memory allocation, heap usage, and garbage collection information.",
			Parameters:  toolParams,
		},
	}
}

func GetServerUptimeTool() api.Tool {
	properties := api.NewToolPropertiesMap()

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_server_uptime",
			Description: "Returns how long the server has been running since it was started.",
			Parameters:  toolParams,
		},
	}
}

func GetConnectedClientsTool() api.Tool {
	properties := api.NewToolPropertiesMap()

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_connected_clients",
			Description: "Returns metadata for all currently connected clients, including their ID, connection time, last activity time, and number of interactions.",
			Parameters:  toolParams,
		},
	}
}

func GetCommandCountTool() api.Tool {
	properties := api.NewToolPropertiesMap()

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_command_count",
			Description: "Returns the total number of commands or requests processed by the server since it started.",
			Parameters:  toolParams,
		},
	}
}
