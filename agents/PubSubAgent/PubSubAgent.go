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

func PubSubAgent(query string, stre *store.MemoryAlloc, user *utils.NewClient, client *api.Client) {
	_ = stre

	ctx := context.Background()
	stream := false

	messages := []api.Message{
		{
			Role: "system",
			Content: `You are an AI pubsub operations assistant for a topic-based messaging server.
Use pubsub tools whenever topic names, subscribers, publishers, topic history, or broadcast actions are needed.
Never invent topic state or client state.
After gathering data, explain the pubsub state clearly and mention operational risks when relevant.`,
		},
		{
			Role:    "user",
			Content: query,
		},
	}

	req := &api.ChatRequest{
		Model:    agents.AGENT_USED,
		Messages: messages,
		Tools: []api.Tool{
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
		},
		Stream:  &stream,
		Options: PubSubInference,
	}

	var toolCallsToExecute []api.ToolCall
	var assistantTextResponse string
	var rawResponses []string

	cerr := client.Chat(ctx, req, func(resp api.ChatResponse) error {
		if len(resp.Message.ToolCalls) > 0 {
			toolCallsToExecute = resp.Message.ToolCalls
		}
		if resp.Message.Content != "" {
			assistantTextResponse = resp.Message.Content
		}
		return nil
	})

	if cerr != nil {
		log.Fatalf("First Ollama API call failed: %v", cerr)
	}

	if len(toolCallsToExecute) > 0 {
		messages = append(messages, api.Message{
			Role:      "assistant",
			ToolCalls: toolCallsToExecute,
		})

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
				toolResult = "No Tools found"
			}

			messages = append(messages, api.Message{
				Role:    "tool",
				Content: toolResult,
			})
			rawResponses = append(rawResponses, toolResult)
		}

		if len(rawResponses) > 0 {
			user.Conn.Write([]byte("\nroxai> " + strings.Join(rawResponses, "\n") + "\n"))
		}

		return
	}

	if assistantTextResponse != "" {
		user.Conn.Write([]byte("\nroxai> " + assistantTextResponse + "\n"))
	}
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
