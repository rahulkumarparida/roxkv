package storageagent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
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

const storageSystemPrompt = `You are the RoxKV storage worker.
You are not a conversational assistant.
Use only the provided storage-analysis tools to inspect persistence, TTL, namespace, and key metadata state.
Return raw tool results only.
Do not summarize, explain, speculate, or invent database facts.
If no available tool can satisfy the request, return a structured error.`

func StorageAgent(query string, stre *store.MemoryAlloc, namespace *store.NameSpace, user *utils.NewClient, client *api.Client) string {
	ctx := context.Background()
	session := agents.GetSession("storageagent", storageSystemPrompt, user)
	tools := []api.Tool{
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
	}

	toolCallsToExecute, assistantTextResponse, cerr := session.Run(ctx, client, agents.AGENT_USED, query, tools, StorageInference)
	if cerr != nil {
		log.Fatalf("Storage worker API call failed: %v", cerr)
	}

	if len(toolCallsToExecute) == 0 {
		if strings.TrimSpace(assistantTextResponse) == "" {
			return agents.StructuredError("no_suitable_tool", "The storage worker could not match the request to an available storage tool.")
		}
		return agents.StructuredError("no_suitable_tool", strings.TrimSpace(assistantTextResponse))
	}

	var rawResults []string
	var results []agents.ToolResult
	for _, tool := range toolCallsToExecute {
		argsByte, _ := json.Marshal(tool.Function.Arguments)
		var toolResult any

		switch tool.Function.Name {
		case "get_total_keys":
			toolResult = metrics.GetTotalKeys(stre, namespace)
		case "get_largest_keys":
			toolResult = metrics.GetLargestKeys(stre, namespace, parseTopN(argsByte))
		case "get_smallest_keys":
			toolResult = metrics.GetSmallestKeys(stre, namespace, parseTopN(argsByte))
		case "get_average_value_size":
			toolResult = metrics.GetAverageValueSize(stre, namespace)
		case "get_namespaces":
			toolResult = metrics.GetNamespaces(stre, namespace)
		case "get_ttl_metrics":
			toolResult = metrics.GetTTLMetrics()
		case "get_expired_keys":
			toolResult = metrics.GetExpiredKeys(stre)
		case "get_upcoming_expirations":
			toolResult = metrics.GetUpcomingExpirations(stre, parseTopN(argsByte))
		case "get_keys_without_ttl":
			toolResult = metrics.GetKeysWithoutTTL(stre)
		case "get_snapshot_count":
			toolResult = metrics.GetSnapshotCount()
		case "get_latest_snapshot":
			toolResult = metrics.GetLatestSnapshot()
		case "get_snapshot_size":
			toolResult = metrics.GetSnapshotSize()
		case "get_persistence_health":
			toolResult = metrics.GetPersistenceHealth()
		case "get_oldest_key":
			toolResult = summarizeKeyMetadata(metrics.GetOldestKey(stre))
		case "get_newest_key":
			toolResult = summarizeKeyMetadata(metrics.GetNewestKey(stre))
		case "get_most_accessed_key":
			toolResult = summarizeKeyMetadata(metrics.GetMostAccessedKey(stre))
		case "get_least_accessed_key":
			toolResult = summarizeKeyMetadata(metrics.GetLeastAccessedKey(stre))
		case "get_recently_modified_keys":
			toolResult = summarizeKeyMetadataSlice(metrics.GetRecentlyModifiedKeys(stre, parseTopN(argsByte)))
		default:
			toolResult = map[string]any{
				"status":  "error",
				"error":   "unknown_tool",
				"message": "The storage worker received an unsupported tool call.",
			}
		}

		rawResults = append(rawResults, marshalStorageResult(toolResult))
		results = append(results, agents.ToolResult{
			Tool:   tool.Function.Name,
			Output: toolResult,
		})
	}

	session.AppendToolResults(rawResults...)
	return agents.StructuredSuccess(results)
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
