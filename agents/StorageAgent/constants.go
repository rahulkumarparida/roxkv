package storageagent

import "github.com/ollama/ollama/api"

var StorageInference = map[string]any{
	"num_predict": 300,
	"num_ctx":     2048,
	"temperature": 0.4,
	"top_p":       0.93,
}

func GetTotalKeysTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_total_keys",
			Description: "Returns the total number of active keys currently available in the database memory.",
			Parameters:  toolParams,
		},
	}
}

func GetLargestKeysTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("top_n", api.ToolProperty{
		Type:        api.PropertyType{"integer"},
		Description: "How many of the largest keys by size to return.",
	})

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
		Required:   []string{"top_n"},
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_largest_keys",
			Description: "Returns the largest keys in memory ranked by recorded value size.",
			Parameters:  toolParams,
		},
	}
}

func GetSmallestKeysTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("top_n", api.ToolProperty{
		Type:        api.PropertyType{"integer"},
		Description: "How many of the smallest keys by size to return.",
	})

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
		Required:   []string{"top_n"},
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_smallest_keys",
			Description: "Returns the smallest keys in memory ranked by recorded value size.",
			Parameters:  toolParams,
		},
	}
}

func GetAverageValueSizeTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_average_value_size",
			Description: "Returns the average recorded value size across active keys in memory.",
			Parameters:  toolParams,
		},
	}
}

func GetNamespacesTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_namespaces",
			Description: "Returns namespace groupings and how many keys belong to each namespace.",
			Parameters:  toolParams,
		},
	}
}

func GetTTLMetricsTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_ttl_metrics",
			Description: "Returns aggregate TTL statistics including active TTL keys, permanent keys, expired counts, and upcoming expirations.",
			Parameters:  toolParams,
		},
	}
}

func GetExpiredKeysTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_expired_keys",
			Description: "Returns the list of keys that are already expired in memory.",
			Parameters:  toolParams,
		},
	}
}

func GetUpcomingExpirationsTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("top_n", api.ToolProperty{
		Type:        api.PropertyType{"integer"},
		Description: "How many upcoming expirations to return.",
	})

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
		Required:   []string{"top_n"},
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_upcoming_expirations",
			Description: "Returns the nearest-expiring keys and their remaining TTL durations.",
			Parameters:  toolParams,
		},
	}
}

func GetKeysWithoutTTLTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_keys_without_ttl",
			Description: "Returns the list of keys that do not have any TTL assigned.",
			Parameters:  toolParams,
		},
	}
}

func GetSnapshotCountTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_snapshot_count",
			Description: "Returns how many snapshot files currently exist on disk.",
			Parameters:  toolParams,
		},
	}
}

func GetLatestSnapshotTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_latest_snapshot",
			Description: "Returns metadata about the most recent snapshot file.",
			Parameters:  toolParams,
		},
	}
}

func GetSnapshotSizeTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_snapshot_size",
			Description: "Returns the total disk size consumed by snapshot files.",
			Parameters:  toolParams,
		},
	}
}

func GetPersistenceHealthTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_persistence_health",
			Description: "Returns persistence health information including folder presence, snapshot count, and latest snapshot details.",
			Parameters:  toolParams,
		},
	}
}

func GetOldestKeyTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_oldest_key",
			Description: "Returns metadata for the oldest key in the active in-memory database.",
			Parameters:  toolParams,
		},
	}
}

func GetNewestKeyTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_newest_key",
			Description: "Returns metadata for the newest key in the active in-memory database.",
			Parameters:  toolParams,
		},
	}
}

func GetMostAccessedKeyTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_most_accessed_key",
			Description: "Returns metadata for the key with the highest recorded access count.",
			Parameters:  toolParams,
		},
	}
}

func GetLeastAccessedKeyTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_least_accessed_key",
			Description: "Returns metadata for the key with the lowest recorded access count.",
			Parameters:  toolParams,
		},
	}
}

func GetRecentlyModifiedKeysTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("top_n", api.ToolProperty{
		Type:        api.PropertyType{"integer"},
		Description: "How many recently modified keys to return.",
	})

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
		Required:   []string{"top_n"},
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_recently_modified_keys",
			Description: "Returns metadata for the most recently modified keys in the active database.",
			Parameters:  toolParams,
		},
	}
}
