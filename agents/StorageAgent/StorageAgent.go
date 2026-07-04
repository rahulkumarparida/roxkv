package storageagent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ollama/ollama/api"
	"github.com/rahulkumarparida/roxkv/agents"
	"github.com/rahulkumarparida/roxkv/internal/metrics"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)

type toolTopNArgs struct {
	TopN int `json:"top_n"`
}

type storageClientSummary struct {
	ID           any    `json:"id"`
	LastUsed     string `json:"last_used"`
	ConnectedAt  string `json:"connected_at"`
	Interactions int    `json:"interactions"`
}

type storageKeyMetadataSummary struct {
	Key            string                `json:"key"`
	TTL            string                `json:"ttl"`
	UpdatedAt      string                `json:"updated_at"`
	CreatedAt      string                `json:"created_at"`
	LastAcessedBy  *storageClientSummary `json:"last_accessed_by,omitempty"`
	KeyAccessCount int64                 `json:"key_access_count"`
	Size           int64                 `json:"size"`
	Namespace      string                `json:"namespace"`
}

func StorageAgent(query string, stre *store.MemoryAlloc, namespace *store.NameSpace, user *utils.NewClient, client *api.Client) {
	ctx := context.Background()
	stream := false

	messages := []api.Message{
		{
			Role: "system",
			Content: `You are an AI storage analysis assistant for an in-memory key-value database.
Use storage tools whenever database storage facts, TTL facts, snapshot facts, or key metadata are needed.
Never estimate or invent database state.
After gathering data, summarize the storage health, important patterns, and any risks or recommendations clearly.`,
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
			GetTotalKeysTool(),
			GetLargestKeysTool(),
			GetSmallestKeysTool(),
			GetAverageValueSizeTool(),
			GetNamespacesTool(),
			GetTTLMetricsTool(),
			GetExpiredKeysTool(),
			GetUpcomingExpirationsTool(),
			GetKeysWithoutTTLTool(),
			GetSnapshotCountTool(),
			GetLatestSnapshotTool(),
			GetSnapshotSizeTool(),
			GetPersistenceHealthTool(),
			GetOldestKeyTool(),
			GetNewestKeyTool(),
			GetMostAccessedKeyTool(),
			GetLeastAccessedKeyTool(),
			GetRecentlyModifiedKeysTool(),
		},
		Stream:  &stream,
		Options: StorageInference,
	}

	var toolCallsToExecute []api.ToolCall
	var assistantTextResponse string

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

		var finalResponse []any

		for _, tool := range toolCallsToExecute {
			argsByte, _ := json.Marshal(tool.Function.Arguments)

			switch tool.Function.Name {
			case "get_total_keys":
				finalResponse = append(finalResponse, metrics.GetTotalKeys(stre, namespace))
			case "get_largest_keys":
				finalResponse = append(finalResponse, metrics.GetLargestKeys(stre, namespace, parseTopN(argsByte)))
			case "get_smallest_keys":
				finalResponse = append(finalResponse, metrics.GetSmallestKeys(stre, namespace, parseTopN(argsByte)))
			case "get_average_value_size":
				finalResponse = append(finalResponse, metrics.GetAverageValueSize(stre, namespace))
			case "get_namespaces":
				finalResponse = append(finalResponse, metrics.GetNamespaces(stre, namespace))
			case "get_ttl_metrics":
				finalResponse = append(finalResponse, metrics.GetTTLMetrics())
			case "get_expired_keys":
				finalResponse = append(finalResponse, metrics.GetExpiredKeys(stre))
			case "get_upcoming_expirations":
				finalResponse = append(finalResponse, metrics.GetUpcomingExpirations(stre, parseTopN(argsByte)))
			case "get_keys_without_ttl":
				finalResponse = append(finalResponse, metrics.GetKeysWithoutTTL(stre))
			case "get_snapshot_count":
				finalResponse = append(finalResponse, metrics.GetSnapshotCount())
			case "get_latest_snapshot":
				finalResponse = append(finalResponse, metrics.GetLatestSnapshot())
			case "get_snapshot_size":
				finalResponse = append(finalResponse, metrics.GetSnapshotSize())
			case "get_persistence_health":
				finalResponse = append(finalResponse, metrics.GetPersistenceHealth())
			case "get_oldest_key":
				finalResponse = append(finalResponse, summarizeKeyMetadata(metrics.GetOldestKey(stre)))
			case "get_newest_key":
				finalResponse = append(finalResponse, summarizeKeyMetadata(metrics.GetNewestKey(stre)))
			case "get_most_accessed_key":
				finalResponse = append(finalResponse, summarizeKeyMetadata(metrics.GetMostAccessedKey(stre)))
			case "get_least_accessed_key":
				finalResponse = append(finalResponse, summarizeKeyMetadata(metrics.GetLeastAccessedKey(stre)))
			case "get_recently_modified_keys":
				finalResponse = append(finalResponse, summarizeKeyMetadataSlice(metrics.GetRecentlyModifiedKeys(stre, parseTopN(argsByte))))
			default:
				continue
			}
		}

		if len(finalResponse) > 0 {
			value := fmt.Sprintf("%v", finalResponse)
			user.Conn.Write([]byte("\nroxai> " + value + "\n"))
		}

		return
	}

	if assistantTextResponse != "" {
		user.Conn.Write([]byte("\nroxai> " + assistantTextResponse + "\n"))
	}
}

func parseTopN(argsByte []byte) int {
	var args toolTopNArgs
	_ = json.Unmarshal(argsByte, &args)
	if args.TopN <= 0 {
		return 5
	}

	return args.TopN
}

func summarizeClient(client *utils.NewClient) *storageClientSummary {
	if client == nil {
		return nil
	}

	return &storageClientSummary{
		ID:           client.ID,
		LastUsed:     client.LastUsed.Format(time.RFC3339),
		ConnectedAt:  client.ConnectedAt.Format(time.RFC3339),
		Interactions: client.Interactions,
	}
}

func summarizeKeyMetadata(meta metrics.KeyMetadata) storageKeyMetadataSummary {
	ttl := ""
	if !meta.TTL.IsZero() {
		ttl = meta.TTL.Format(time.RFC3339)
	}

	createdAt := ""
	if !meta.CreatedAt.IsZero() {
		createdAt = meta.CreatedAt.Format(time.RFC3339)
	}

	updatedAt := ""
	if !meta.UpdatedAt.IsZero() {
		updatedAt = meta.UpdatedAt.Format(time.RFC3339)
	}

	return storageKeyMetadataSummary{
		Key:            meta.Key,
		TTL:            ttl,
		UpdatedAt:      updatedAt,
		CreatedAt:      createdAt,
		LastAcessedBy:  summarizeClient(meta.LastAcessedBy),
		KeyAccessCount: meta.KeyAccessCount,
		Size:           meta.Size,
		Namespace:      meta.Namespace,
	}
}

func summarizeKeyMetadataSlice(items []metrics.KeyMetadata) []storageKeyMetadataSummary {
	summaries := make([]storageKeyMetadataSummary, 0, len(items))
	for _, item := range items {
		summaries = append(summaries, summarizeKeyMetadata(item))
	}

	return summaries
}

func marshalStorageResult(value any) string {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", value)
	}

	return string(data)
}
