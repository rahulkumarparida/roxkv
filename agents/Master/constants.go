package master

import "github.com/ollama/ollama/api"

var MasterInference = map[string]any{
	"num_predict": 400,
	"num_ctx":     4096,
	"temperature": 0.4,
	"top_p":       0.9,
}

func KvAgentTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("query", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Use this capability when the master agent must execute direct key-value operations such as reading a specific value, storing a new value, deleting a key, listing all keys, saving the current store, or restoring previously saved data. Return the raw KV result so the master agent can present a concise answer or a deeper analysis.",
	})

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
		Required:   []string{"query"},
	}

	return api.Tool{
		Type: "agent",
		Function: api.ToolFunction{
			Name:        "kvagent",
			Description: "Purpose: delegate direct key-value execution. Inputs: query string describing a read, write, delete, list, save, or load action. Output: structured raw execution data from the KV worker. Use when the request is about exact key operations or persistence commands. Do not use for monitoring, pubsub, or storage analytics.",
			Parameters:  toolParams,
		},
	}
}

func MonitorAgentTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("query", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Use this capability when the master agent needs live system or runtime information such as CPU load, memory usage, disk usage, uptime, client activity, or general server health. Return the raw monitoring result so the master agent can explain the state clearly.",
	})

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
		Required:   []string{"query"},
	}

	return api.Tool{
		Type: "agent",
		Function: api.ToolFunction{
			Name:        "monitoragent",
			Description: "Purpose: delegate live monitoring execution. Inputs: query string describing the exact system or runtime metric needed. Output: structured raw execution data from the monitoring worker. Use when the request is about CPU, RAM, disk, runtime stats, uptime, clients, or command counts. Do not use for key mutation or pubsub operations.",
			Parameters:  toolParams,
		},
	}
}

func PubSubAgentTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("query", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Use this capability when the master agent must inspect or change topic-based messaging state, including topics, publishers, subscribers, history, broadcasts, topic creation, topic deletion, or user removal. Return the raw pubsub result so the master agent can summarize the messaging outcome.",
	})

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
		Required:   []string{"query"},
	}

	return api.Tool{
		Type: "agent",
		Function: api.ToolFunction{
			Name:        "pubsubagent",
			Description: "Purpose: delegate pubsub execution. Inputs: query string describing the exact topic, subscriber, publisher, history, broadcast, or topic-lifecycle action. Output: structured raw execution data from the pubsub worker. Use when the request is about topic-based messaging state or actions. Do not use for key-value data or server metrics.",
			Parameters:  toolParams,
		},
	}
}

func StorageAgentTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("query", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Use this capability when the master agent needs persistence, storage, TTL, namespace, snapshot, or key-metadata analysis rather than simple get/set/delete operations. Return the raw storage result so the master agent can explain storage health and risks.",
	})

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
		Required:   []string{"query"},
	}

	return api.Tool{
		Type: "agent",
		Function: api.ToolFunction{
			Name:        "storageagent",
			Description: "Purpose: delegate storage and persistence analysis. Inputs: query string describing TTL, namespace, snapshot, size, or key-metadata inspection. Output: structured raw execution data from the storage worker. Use when the request is about storage health or persistence facts rather than direct key mutation. Do not use for simple get/set/delete requests or live runtime metrics.",
			Parameters:  toolParams,
		},
	}
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
			Description: "Master agent capability: provide a high-level summary of the current host, including the username, operating system, architecture, CPU count, memory state, and runtime context. Use this for broad machine information questions.",
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
			Description: "Master agent capability: provide the current CPU utilization percentage. Use this when the request is specifically about CPU load rather than overall machine health.",
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
			Description: "Master agent capability: provide the current RAM usage details, including total memory, free memory, and the usage percentage. Use this for memory-focused questions.",
			Parameters:  toolParams,
		},
	}
}

func GetDiskUsageTool() api.Tool {
	properties := api.NewToolPropertiesMap()

	pathProp := api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Master agent capability: inspect disk capacity and availability for one filesystem path such as / or /home. Use this when the request asks about disk usage or storage availability.",
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
			Description: "Master agent capability: return Go runtime statistics for the RoxKV process, including version, OS, architecture, goroutine count, allocation, heap, and garbage collection details. Use this for runtime/internal process questions.",
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
			Description: "Master agent capability: report how long the RoxKV server has been running. Use this when uptime is the requested metric.",
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
			Description: "Master agent capability: report connected client metadata including IDs, connection time, last activity, and interaction counts. Use this for active-client or session monitoring questions.",
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
			Description: "Master agent capability: report the total number of commands or requests processed by the server since startup. Use this for activity-volume questions.",
			Parameters:  toolParams,
		},
	}
}




func GetAllTopicsOnlineTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_all_topics_online",
			Description: "Master agent capability: list every currently available pubsub topic. Use this when the request asks which topics exist.",
			Parameters:  toolParams,
		},
	}
}

func GetAllSubscribersTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_all_subscribers",
			Description: "Master agent capability: return subscriber information across all topics. Use this for global subscriber inspection.",
			Parameters:  toolParams,
		},
	}
}

func GetAllPublishersTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_all_publishers",
			Description: "Master agent capability: return publisher information across all topics. Use this for global publisher inspection.",
			Parameters:  toolParams,
		},
	}
}

func GetPublishersTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("topic", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Master agent capability: return the publishers for one specific topic. Use this when the request names a single topic.",
	})

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
		Required:   []string{"topic"},
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_publishers",
			Description: "Purpose: return publishers for one topic. Inputs: topic string. Output: publisher client metadata for that topic. Use when the request names a single topic. Do not use for all-topic listings or message history.",
			Parameters:  toolParams,
		},
	}
}

func GetSubscribersTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("topic", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Master agent capability: return the subscribers for one specific topic. Use this when the request names a single topic.",
	})

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
		Required:   []string{"topic"},
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_subscribers",
			Description: "Purpose: return subscribers for one topic. Inputs: topic string. Output: subscriber client metadata for that topic. Use when the request names a single topic. Do not use for global subscriber listings or publisher data.",
			Parameters:  toolParams,
		},
	}
}

func GetTopicStatisticsTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("topic", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Master agent capability: return a complete snapshot of a topic, including publishers, subscribers, message history, timestamps, and size. Use this for detailed topic analysis.",
	})

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
		Required:   []string{"topic"},
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_topic_statistics",
			Description: "Purpose: return a full topic snapshot. Inputs: topic string. Output: publishers, subscribers, publish count, history, timestamps, and total message size for that topic. Use when the request needs detailed state for one topic. Do not use for simple existence checks.",
			Parameters:  toolParams,
		},
	}
}

func GetInactiveTopicsTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "get_inactive_topics",
			Description: "Master agent capability: identify older or idle topics that may be stale. Use this when the request asks for inactive topics.",
			Parameters:  toolParams,
		},
	}
}

func BroadcastToTopicTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("topic", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Master agent capability: publish one message to a single topic. Use this for targeted broadcasts.",
	})
	properties.Set("message", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Message payload to publish to that topic.",
	})

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
		Required:   []string{"topic", "message"},
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "broadcast_to_topic",
			Description: "Purpose: publish one message to one topic. Inputs: topic string, message string. Output: broadcast result from the pubsub layer. Use when the request explicitly asks to send to a single topic. Do not use for all-topic broadcast or read-only inspection.",
			Parameters:  toolParams,
		},
	}
}

func BroadcastEverywhereTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("message", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Master agent capability: publish one message across all topics. Use this for global broadcasts.",
	})

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
		Required:   []string{"message"},
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "broadcast_everywhere",
			Description: "Purpose: publish one message across all topics. Inputs: message string. Output: broadcast result from the pubsub layer. Use when the request explicitly asks for a global broadcast. Do not use for one-topic sends or topic inspection.",
			Parameters:  toolParams,
		},
	}
}

func HistoryOfTopicTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("topic", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Master agent capability: return the recorded publish history for one topic. Use this when the request asks what was sent on a topic.",
	})

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
		Required:   []string{"topic"},
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "history_of_topic",
			Description: "Purpose: return publish history for one topic. Inputs: topic string. Output: recorded messages with sizes and timestamps. Use when the request asks what was published on a topic. Do not use for subscriber or publisher membership checks.",
			Parameters:  toolParams,
		},
	}
}

func DeleteTopicTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("topic", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Master agent capability: delete one pubsub topic. Use this when the request explicitly asks to remove a topic.",
	})

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
		Required:   []string{"topic"},
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "delete_topic",
			Description: "Purpose: delete one pubsub topic. Inputs: topic string. Output: deletion result from the pubsub layer. Use when the request explicitly asks to remove a topic. Do not use for message cleanup, history reads, or topic creation.",
			Parameters:  toolParams,
		},
	}
}

func CreateTopicTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("topic", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Master agent capability: create one new pubsub topic. Use this when the request explicitly asks to create a topic.",
	})

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
		Required:   []string{"topic"},
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "create_topic",
			Description: "Purpose: create one new pubsub topic. Inputs: topic string. Output: topic snapshot after creation. Use when the request explicitly asks to create a topic. Do not use for inspection, deletion, or broadcasting.",
			Parameters:  toolParams,
		},
	}
}

func RemoveUserTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("portaddress", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Master agent capability: disconnect one user or client from the pubsub layer and optionally explain the reason. Use this when the request asks to remove a user.",
	})
	properties.Set("reason", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Reason message to send before disconnecting the user.",
	})

	toolParams := api.ToolFunctionParameters{
		Type:       "object",
		Properties: properties,
		Required:   []string{"portaddress", "reason"},
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "remove_user",
			Description: "Purpose: disconnect one user from the pubsub layer. Inputs: portaddress string, reason string. Output: removal result from the pubsub layer. Use when the request explicitly asks to remove a client. Do not use for unsubscribing by topic or simple inspection.",
			Parameters:  toolParams,
		},
	}
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
			Description: "Master agent capability: count the active keys currently stored in memory. Use this when the request asks for a total key count.",
			Parameters:  toolParams,
		},
	}
}

func GetLargestKeysTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("top_n", api.ToolProperty{
		Type:        api.PropertyType{"integer"},
		Description: "Master agent capability: return the largest keys by stored value size. Use this when the request asks which keys consume the most space.",
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
		Description: "Master agent capability: return the smallest keys by stored value size. Use this when the request asks which keys are smallest.",
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
			Description: "Master agent capability: calculate the average size of stored values. Use this when the request asks about typical value size.",
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
			Description: "Master agent capability: summarize key distribution by namespace. Use this when the request asks how data is grouped.",
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
			Description: "Master agent capability: report TTL-related metrics such as active TTL keys, permanent keys, expired keys, and upcoming expirations. Use this for TTL health questions.",
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
			Description: "Master agent capability: list keys whose TTL has already expired. Use this when the request asks which keys are expired.",
			Parameters:  toolParams,
		},
	}
}

func GetUpcomingExpirationsTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("top_n", api.ToolProperty{
		Type:        api.PropertyType{"integer"},
		Description: "Master agent capability: list keys that are about to expire soon. Use this when the request asks what will expire next.",
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
			Description: "Master agent capability: list keys that have no TTL assigned. Use this when the request asks which keys are permanent or unbounded.",
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
			Description: "Master agent capability: count snapshot files stored on disk. Use this when the request asks how many snapshots exist.",
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
			Description: "Master agent capability: return the most recent snapshot metadata. Use this when the request asks for the latest snapshot.",
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
			Description: "Master agent capability: report the total disk space consumed by snapshots. Use this when the request asks about snapshot storage usage.",
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
			Description: "Master agent capability: summarize persistence health, including folder presence, snapshot count, and latest snapshot status. Use this for persistence health checks.",
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
			Description: "Master agent capability: return metadata for the oldest active key. Use this when the request asks which key is oldest.",
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
			Description: "Master agent capability: return metadata for the newest active key. Use this when the request asks which key is newest.",
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
			Description: "Master agent capability: return metadata for the most frequently accessed key. Use this when the request asks which key is most used.",
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
			Description: "Master agent capability: return metadata for the least frequently accessed key. Use this when the request asks which key is least used.",
			Parameters:  toolParams,
		},
	}
}

func GetRecentlyModifiedKeysTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("top_n", api.ToolProperty{
		Type:        api.PropertyType{"integer"},
		Description: "Master agent capability: list keys that were modified most recently. Use this when the request asks what changed recently.",
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



func GetKeyTool() api.Tool {

	var properties = api.NewToolPropertiesMap()

	keyprop := api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Master agent capability: read one exact value from the in-memory store by key. Use this when the request asks for a specific stored value.",
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
		Description: "Master agent capability: create or overwrite one key-value pair in memory. Use this when the request asks to save, set, or update a value.",
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
		Description: "Master agent capability: delete one key and its stored value from memory. Use this when the request explicitly asks to remove a key.",
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
			Description: "Master agent capability: list all keys currently present in memory. Use this when the request asks to enumerate stored keys.",
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
			Description: "Master agent capability: persist the current in-memory database state to disk. Use this when the request asks to save the current data.",
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
			Description: "Master agent capability: restore previously saved database data from disk back into memory. Use this when the request asks to load or recover persisted state.",
			Parameters:  toolParams,
		},
	}

	return LoadDiskTool
}




