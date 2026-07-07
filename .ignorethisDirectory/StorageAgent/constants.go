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
			Description: "Purpose: count active keys in memory. Inputs: none. Output: integer key count. Use when the request asks for total keys. Do not use for listing keys or key-size rankings.",
			Parameters:  toolParams,
		},
	}
}

func GetLargestKeysTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("top_n", api.ToolProperty{
		Type:        api.PropertyType{"integer"},
		Description: "Positive integer count of largest keys to return.",
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
			Description: "Purpose: rank the largest keys by stored value size. Inputs: top_n integer. Output: top key-size records in descending size order. Use when the request asks which keys consume the most memory. Do not use for smallest keys or total counts.",
			Parameters:  toolParams,
		},
	}
}

func GetSmallestKeysTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("top_n", api.ToolProperty{
		Type:        api.PropertyType{"integer"},
		Description: "Positive integer count of smallest keys to return.",
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
			Description: "Purpose: rank the smallest keys by stored value size. Inputs: top_n integer. Output: top key-size records in ascending size order. Use when the request asks which keys are smallest. Do not use for largest keys or total counts.",
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
			Description: "Purpose: compute average value size across active keys. Inputs: none. Output: average size as an integer. Use when the request asks for typical value size. Do not use for per-key rankings.",
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
			Description: "Purpose: summarize key distribution by namespace. Inputs: none. Output: namespace names with key counts. Use when the request asks how keys are grouped. Do not use for single-key metadata.",
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
			Description: "Purpose: return aggregate TTL metrics. Inputs: none. Output: active TTL count, permanent key count, expired count, and upcoming expiration totals. Use when the request is about TTL health. Do not use for listing exact keys.",
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
			Description: "Purpose: list keys whose TTL has already expired. Inputs: none. Output: array of expired key names. Use when the request asks which keys are expired. Do not use for future expirations.",
			Parameters:  toolParams,
		},
	}
}

func GetUpcomingExpirationsTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("top_n", api.ToolProperty{
		Type:        api.PropertyType{"integer"},
		Description: "Positive integer count of upcoming expirations to return.",
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
			Description: "Purpose: list the nearest upcoming expirations. Inputs: top_n integer. Output: key names with remaining TTL duration. Use when the request asks what will expire soon. Do not use for already expired keys or lifetime totals.",
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
			Description: "Purpose: list keys with no TTL. Inputs: none. Output: array of key names without expiration. Use when the request asks which keys are permanent. Do not use for TTL counts only.",
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
			Description: "Purpose: count snapshot files on disk. Inputs: none. Output: integer snapshot count. Use when the request asks how many persistence snapshots exist. Do not use for snapshot size or latest snapshot details.",
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
			Description: "Purpose: return metadata for the newest snapshot file. Inputs: none. Output: latest snapshot details from disk. Use when the request asks for the most recent snapshot. Do not use for full persistence health or snapshot totals.",
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
			Description: "Purpose: return total disk space consumed by snapshots. Inputs: none. Output: aggregate snapshot size. Use when the request asks about snapshot storage usage. Do not use for key-size analysis.",
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
			Description: "Purpose: return persistence health summary. Inputs: none. Output: snapshot-folder presence, snapshot count, and latest snapshot details. Use when the request asks whether persistence looks healthy. Do not use for per-key TTL analysis.",
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
			Description: "Purpose: return metadata for the oldest active key. Inputs: none. Output: key lifecycle metadata. Use when the request asks for the oldest entry. Do not use for recently modified or access-frequency questions.",
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
			Description: "Purpose: return metadata for the newest active key. Inputs: none. Output: key lifecycle metadata. Use when the request asks for the newest entry. Do not use for oldest or access-frequency questions.",
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
			Description: "Purpose: return metadata for the most accessed key. Inputs: none. Output: key lifecycle and access metadata. Use when the request asks which key is hottest. Do not use for creation-time rankings.",
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
			Description: "Purpose: return metadata for the least accessed key. Inputs: none. Output: key lifecycle and access metadata. Use when the request asks which key is least used. Do not use for creation-time rankings.",
			Parameters:  toolParams,
		},
	}
}

func GetRecentlyModifiedKeysTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("top_n", api.ToolProperty{
		Type:        api.PropertyType{"integer"},
		Description: "Positive integer count of recently modified keys to return.",
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
			Description: "Purpose: list recently modified keys. Inputs: top_n integer. Output: key metadata ordered by recent modification time. Use when the request asks what changed most recently. Do not use for oldest/newest creation-time comparisons.",
			Parameters:  toolParams,
		},
	}
}
