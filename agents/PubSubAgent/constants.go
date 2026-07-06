package pubsubagent

import "github.com/ollama/ollama/api"

var PubSubInference = map[string]any{
	"num_predict": 200,
	"num_ctx":     2048,
	"temperature": 0.4,
	"top_p":       0.92,
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
			Description: "Purpose: list all topics currently present in the pubsub system. Inputs: none. Output: array of topic names. Use when the request asks which topics exist. Do not use for subscriber details or message history.",
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
			Description: "Purpose: return subscriber snapshots across all topics. Inputs: none. Output: subscriber client metadata grouped by topic. Use when the request asks for all subscribers globally. Do not use for one topic only or publisher data.",
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
			Description: "Purpose: return publisher snapshots across all topics. Inputs: none. Output: publisher client metadata grouped by topic. Use when the request asks for all publishers globally. Do not use for one topic only or subscriber data.",
			Parameters:  toolParams,
		},
	}
}

func GetPublishersTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("topic", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Exact topic name whose publishers should be returned.",
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
		Description: "Exact topic name whose subscribers should be returned.",
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
		Description: "Exact topic name whose full statistics should be returned.",
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
			Description: "Purpose: return topic snapshots ordered for inactivity inspection. Inputs: none. Output: topic snapshots suitable for finding old or idle topics. Use when the request asks for inactive or stale topics. Do not use for active broadcast actions.",
			Parameters:  toolParams,
		},
	}
}

func BroadcastToTopicTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("topic", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Exact topic name that should receive the message.",
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
		Description: "Message payload to publish across every topic.",
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
		Description: "Exact topic name whose publish history should be returned.",
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
		Description: "Exact topic name to delete from the pubsub system.",
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
		Description: "Exact topic name to create.",
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
		Description: "Client ID or port address of the user to disconnect.",
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
