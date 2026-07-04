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
			Description: "Returns the names of all currently available pubsub topics that exist on the server.",
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
			Description: "Returns the current subscriber client data across all topics, including real client field values instead of pointer addresses.",
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
			Description: "Returns the current publisher client data across all topics, including real client field values instead of pointer addresses.",
			Parameters:  toolParams,
		},
	}
}

func GetPublishersTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("topic", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Topic name whose publishers should be returned.",
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
			Description: "Returns publisher client data for a specific topic.",
			Parameters:  toolParams,
		},
	}
}

func GetSubscribersTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("topic", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Topic name whose subscribers should be returned.",
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
			Description: "Returns subscriber client data for a specific topic.",
			Parameters:  toolParams,
		},
	}
}

func GetTopicStatisticsTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("topic", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Topic name whose pubsub statistics should be returned.",
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
			Description: "Returns a full snapshot of a topic including subscribers, publishers, publish count, history, timestamps, and total message size.",
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
			Description: "Returns topic snapshots ordered by creation time so older and likely inactive topics can be inspected.",
			Parameters:  toolParams,
		},
	}
}

func BroadcastToTopicTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("topic", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Topic name that should receive the broadcast message.",
	})
	properties.Set("message", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Message content to publish to the topic.",
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
			Description: "Publishes a message to one specific topic and returns whether the broadcast request succeeded.",
			Parameters:  toolParams,
		},
	}
}

func BroadcastEverywhereTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("message", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Message content to publish across all topics.",
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
			Description: "Publishes one message across every topic on the server and returns whether the broadcast request succeeded.",
			Parameters:  toolParams,
		},
	}
}

func HistoryOfTopicTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("topic", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Topic name whose message history should be returned.",
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
			Description: "Returns the recorded publish history for a topic including message contents, sizes, and timestamps.",
			Parameters:  toolParams,
		},
	}
}

func DeleteTopicTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("topic", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Topic name to close and remove from the server.",
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
			Description: "Deletes or closes a pubsub topic from the server and returns whether the operation succeeded.",
			Parameters:  toolParams,
		},
	}
}

func CreateTopicTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("topic", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Topic name to create.",
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
			Description: "Creates a new pubsub topic and returns the topic snapshot after creation.",
			Parameters:  toolParams,
		},
	}
}

func RemoveUserTool() api.Tool {
	properties := api.NewToolPropertiesMap()
	properties.Set("portaddress", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Client ID or port address of the user that should be disconnected.",
	})
	properties.Set("reason", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Reason message that will be sent to the user before closing the connection.",
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
			Description: "Disconnects a subscriber by client ID or port address and returns whether the removal request succeeded.",
			Parameters:  toolParams,
		},
	}
}
