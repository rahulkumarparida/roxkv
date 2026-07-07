package pubsubagent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/rahulkumarparida/roxkv/agents"
	"github.com/rahulkumarparida/roxkv/internal/metrics"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)

type pubsubTopicArgs struct {
	Topic string `json:"topic"`
}

type pubsubMessageArgs struct {
	Message string `json:"message"`
}

type pubsubTopicMessageArgs struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

type pubsubRemoveUserArgs struct {
	PortAddress string `json:"portaddress"`
	Reason      string `json:"reason"`
}

const pubsubSystemPrompt = `You are the RoxKV pubsub worker.
You are not a conversational assistant.
Use only the provided pubsub tools to inspect or mutate topic state.
Return raw tool results only.
Do not summarize, explain, speculate, or invent topic data.
If no available tool can satisfy the request, return a structured error.`

func PubSubAgent(query string, stre *store.MemoryAlloc, user *utils.NewClient, client *api.Client) string {
	_ = stre

	ctx := context.Background()
	session := agents.GetSession("pubsubagent", pubsubSystemPrompt, user)
	tools := []api.Tool{
		GetAllTopicsOnlineTool(),
		GetAllSubscribersTool(),
		GetAllPublishersTool(),
		GetPublishersTool(),
		GetSubscribersTool(),
		GetTopicStatisticsTool(),
		GetInactiveTopicsTool(),
		BroadcastToTopicTool(),
		BroadcastEverywhereTool(),
		HistoryOfTopicTool(),
		DeleteTopicTool(),
		CreateTopicTool(),
		RemoveUserTool(),
	}

	toolCallsToExecute, assistantTextResponse, cerr := session.Run(ctx, client, agents.AGENT_USED, query, tools, PubSubInference)
	if cerr != nil {
		log.Fatalf("PubSub worker API call failed: %v", cerr)
	}

	if len(toolCallsToExecute) == 0 {
		if strings.TrimSpace(assistantTextResponse) == "" {
			return agents.StructuredError("no_suitable_tool", "The pubsub worker could not match the request to an available pubsub tool.")
		}
		return agents.StructuredError("no_suitable_tool", strings.TrimSpace(assistantTextResponse))
	}

	var rawResults []string
	var results []agents.ToolResult
	for _, tool := range toolCallsToExecute {
		argsByte, _ := json.Marshal(tool.Function.Arguments)
		var toolResult string

		switch tool.Function.Name {
		case "get_all_topics_online":
			toolResult = marshalPubSubResult(metrics.GetAllTopicsOnline(user))
		case "get_all_subscribers":
			toolResult = marshalPubSubResult(metrics.GetAllSubscribers(user))
		case "get_all_publishers":
			toolResult = marshalPubSubResult(metrics.GetAllPublisher(user))
		case "get_publishers":
			toolResult = marshalPubSubResult(metrics.GetPublishers(user, parseTopic(argsByte)))
		case "get_subscribers":
			toolResult = marshalPubSubResult(metrics.GetSubscribers(user, parseTopic(argsByte)))
		case "get_topic_statistics":
			toolResult = marshalPubSubResult(metrics.GetTopicStatistics(user, parseTopic(argsByte)))
		case "get_inactive_topics":
			toolResult = marshalPubSubResult(metrics.GetInactiveTopics(user))
		case "broadcast_to_topic":
			topic, message := parseTopicMessage(argsByte)
			toolResult = marshalPubSubResult(metrics.BroadcastToTopic(user, topic, message))
		case "broadcast_everywhere":
			toolResult = marshalPubSubResult(metrics.BroadcastEverywhere(user, parseMessage(argsByte)))
		case "history_of_topic":
			toolResult = marshalPubSubResult(metrics.HistoryOfTopic(user, parseTopic(argsByte)))
		case "delete_topic":
			toolResult = marshalPubSubResult(metrics.DeleteTopic(user, parseTopic(argsByte)))
		case "create_topic":
			toolResult = marshalPubSubResult(metrics.CreateTopic(user, parseTopic(argsByte)))
		case "remove_user":
			portAddress, reason := parseRemoveUserArgs(argsByte)
			toolResult = marshalPubSubResult(metrics.RemoveUser(user, portAddress, reason))
		default:
			toolResult = agents.StructuredError("unknown_tool", "The pubsub worker received an unsupported tool call.")
		}

		rawResults = append(rawResults, toolResult)
		results = append(results, agents.ToolResult{
			Tool:   tool.Function.Name,
			Output: agents.ParseToolOutput(toolResult),
		})
	}

	session.AppendToolResults(rawResults...)
	return agents.StructuredSuccess(results)
}

func parseTopic(argsByte []byte) string {
	var args pubsubTopicArgs
	_ = json.Unmarshal(argsByte, &args)
	return strings.TrimSpace(args.Topic)
}

func parseMessage(argsByte []byte) string {
	var args pubsubMessageArgs
	_ = json.Unmarshal(argsByte, &args)
	return strings.TrimSpace(args.Message)
}

func parseTopicMessage(argsByte []byte) (string, string) {
	var args pubsubTopicMessageArgs
	_ = json.Unmarshal(argsByte, &args)
	return strings.TrimSpace(args.Topic), strings.TrimSpace(args.Message)
}

func parseRemoveUserArgs(argsByte []byte) (string, string) {
	var args pubsubRemoveUserArgs
	_ = json.Unmarshal(argsByte, &args)
	return strings.TrimSpace(args.PortAddress), strings.TrimSpace(args.Reason)
}

func marshalPubSubResult(value any) string {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", value)
	}

	return string(data)
}
