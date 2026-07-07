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
			Description: "Purpose: return a host summary for the current machine. Inputs: none. Output: username, OS, architecture, CPU count, RAM stats, and goroutine count. Use when the request asks for overall machine information. Do not use for focused CPU, RAM, disk, or uptime queries.",
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
			Description: "Purpose: return current CPU utilization. Inputs: none. Output: CPU usage percentage as a string. Use when the request is specifically about CPU load. Do not use for general host summaries or memory metrics.",
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
			Description: "Purpose: return current RAM statistics. Inputs: none. Output: total RAM, free RAM, and used percentage. Use when the request is specifically about memory usage. Do not use for disk, CPU, or uptime questions.",
			Parameters:  toolParams,
		},
	}
}

func GetDiskUsageTool() api.Tool {
	properties := api.NewToolPropertiesMap()

	pathProp := api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Filesystem path whose disk stats should be inspected, such as '/' or '/home'.",
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
			Description: "Purpose: return disk capacity details for one filesystem path. Inputs: path string. Output: total, free, and available space in GB for that path. Use when the request is about disk usage for a specific mount or directory path. Do not use for RAM, CPU, or runtime stats.",
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
			Description: "Purpose: return Go runtime statistics for the RoxKV process. Inputs: none. Output: Go version, OS, arch, CPU count, goroutines, allocation, heap, and GC counters. Use when the request is about process/runtime internals. Do not use for host disk usage or client lists.",
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
			Description: "Purpose: return RoxKV server uptime. Inputs: none. Output: elapsed duration since startup. Use when the request asks how long the server has been running. Do not use for resource metrics or client metadata.",
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
			Description: "Purpose: return metadata for all connected clients. Inputs: none. Output: client IDs, connection times, last-used timestamps, and interaction counts. Use when the request is about active chat/server clients. Do not use for pubsub subscriber lists or system load.",
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
			Description: "Purpose: return the total number of requests processed since startup. Inputs: none. Output: integer command count. Use when the request asks for overall server activity volume. Do not use for per-client activity or key counts.",
			Parameters:  toolParams,
		},
	}
}
