package master

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ollama/ollama/api"
	"github.com/rahulkumarparida/roxkv/agents"
	"github.com/rahulkumarparida/roxkv/agents/registry"
	"github.com/rahulkumarparida/roxkv/agents/router"
	"github.com/rahulkumarparida/roxkv/internal/commands"
	"github.com/rahulkumarparida/roxkv/internal/logger"
	"github.com/rahulkumarparida/roxkv/internal/metrics"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)

type masterToolQueryArgs struct {
	Query string `json:"query"`
}

// Pubsub
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

// Storage
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

const masterSystemPrompt = `You are RoxAI, the intelligent orchestration layer of RoxKV.

RoxKV is a custom in-memory database written in Go that provides key-value storage, persistence, TTL management, monitoring, Pub/Sub messaging, snapshots, and system diagnostics through executable tools.

Your responsibility is not to perform these operations yourself. Your responsibility is to understand the user's intent, decide which tools are required, execute them, analyze the returned results, and provide a clear response.

Workflow:

1. Understand the user's request.
2. Determine which tool or combination of tools is required.
3. Execute only the necessary tools.
4. Wait for the tool results.
5. Analyze and correlate the returned data.
6. Respond naturally with a single, coherent answer.

Guidelines:

- Always rely on tool output instead of assumptions.
- Never invent data or system state.
- Never claim an operation succeeded unless the tool confirms it.
- If multiple tools are required, combine their results into one response.
- If a request is ambiguous, ask for clarification before executing tools.
- Keep responses concise unless the user requests a detailed explanation.
- Explain errors honestly and suggest the next step when appropriate.

Summarization Rules:

The results returned by tools are the source of truth.

After receiving tool outputs:

1. Read and understand all returned data.
2. Correlate information across multiple tool results when applicable.
3. Summarize the information into a natural, easy-to-read response.
4. Highlight only the information relevant to the user's request.
5. Avoid repeating raw data unless the user explicitly asks for it.
6. When useful, provide brief observations, insights, or recommendations based only on the returned data.
7. If the user requests a detailed report, include all important findings in a structured format.
8. If the user requests a short answer, provide only the essential information.
9. Never modify, exaggerate, or invent values that are not present in the tool results.

Treat tool outputs as factual evidence. Your job is to transform raw structured data into clear, human-readable information while preserving accuracy.

You are aware that every tool represents a capability of RoxKV. Your job is to coordinate these capabilities, not replace them.

Your goal is to help users interact with RoxKV naturally, accurately, and safely.`

func MasterAgent(query string, stre *store.MemoryAlloc, user *utils.NewClient, client *api.Client) any {
	logger.InfoLog("Master Agent called with query: " + query)
	ctx := context.Background()
	session := agents.GetSession("masteragent", masterSystemPrompt, user)
	// Route tools: score every registered tool against the user query
	// and expose only the relevant subset to the LLM.
	tools := router.RouteTools(query, registry.AllMetadata(), router.DefaultMaxTools)
	var finalizedTools []api.Tool
	if len(tools) > 6 {

		finalizedTools = tools[1:6]
	} else {
		finalizedTools = tools
	}

	fmt.Println("Sending it to the session runner with total ", len(tools), " tools-->")
	toolsToExecute, assistantTextResponse, cerr := session.Run(ctx, client, agents.AGENT_USED, query, finalizedTools, MasterInference)

	if cerr != nil {
		log.Fatalf("Master orchestrator API call failed: %v", cerr)
	}

	if len(toolsToExecute) > 0 {
		var rawReplies []string
		var specialistReplies []string
		for _, tool := range toolsToExecute {
			fmt.Println("Executing: ", tool)

			argsByte, _ := json.Marshal(tool.Function.Arguments)

			var args masterToolQueryArgs
			_ = json.Unmarshal(argsByte, &args)
			fmt.Println("Arguments Here :", args)

			queryText := strings.TrimSpace(args.Query)

			if queryText == "" {
				queryText = "Please inspect the relevant state for this request."
			}

			var response any

			switch tool.Function.Name {
			case "get":
				var args struct {
					Key string `json:"key"`
				}
				_ = json.Unmarshal(argsByte, &args)
				result := store.GetKv(stre, strings.TrimSpace(args.Key))
				response = result.Val
			case "set":
				var args struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				}
				_ = json.Unmarshal(argsByte, &args)

				meta := store.Metadata{
					TTL:           time.Time{},
					UpdatedAt:     time.Now(),
					LastAcessedBy: user,
					Size:          int64(len(args.Value)),
				}
				item := store.Item{
					Key:  strings.TrimSpace(args.Key),
					Val:  strings.TrimSpace(args.Value),
					Meta: meta,
				}

				store.TTLMetricsContainer.PermanentKeys += 1
				response = strconv.FormatBool(store.SetKv(stre, &item))
			case "del":
				var args struct {
					Key string `json:"key"`
				}
				_ = json.Unmarshal(argsByte, &args)
				response = strconv.FormatBool(store.DelKv(stre, strings.TrimSpace(args.Key)))
			case "keys":
				response = store.KeyKv(stre)
			case "save":
				response = commands.SaveCommand(stre)
			case "load":

				response = commands.LoaderCommand(stre)
			case "get_computer_usage":
				response = marshalMonitorResult(metrics.GetComputerUsage())
			case "get_cpu_usage":
				response = metrics.GetCPUUsage()
			case "get_ram_usage":
				response = marshalMonitorResult(metrics.GetRAMUsage())
			case "get_disk_usage":
				var args struct {
					Path string `json:"path"`
				}
				_ = json.Unmarshal(argsByte, &args)

				path := strings.TrimSpace(args.Path)
				if path == "" {
					path = "/"
				}

				disk := metrics.GetDiskUsage(path)
				if disk.Err != nil {
					response = agents.StructuredError("disk_usage_failed", fmt.Sprintf("failed to fetch disk usage for %q: %v", path, disk.Err))
				} else {
					response = marshalMonitorResult(map[string]any{
						"path":         path,
						"total_gb":     disk.Total,
						"free_gb":      disk.Free,
						"available_gb": disk.Avaliable,
					})
				}
			case "get_runtime_stats":
				response = marshalMonitorResult(metrics.GetRuntimeStats())
			case "get_server_uptime":
				response = fmt.Sprintf("%.2f", metrics.GetServerUptime())
			case "get_connected_clients":
				response = marshalMonitorResult(metrics.GetConnectedClients())
			case "get_command_count":
				response = fmt.Sprintf("%d", metrics.GetCommandCount())
			case "get_all_topics_online":
				response = marshalPubSubResult(metrics.GetAllTopicsOnline(user))
			case "get_all_subscribers":
				response = marshalPubSubResult(metrics.GetAllSubscribers(user))
			case "get_all_publishers":
				response = marshalPubSubResult(metrics.GetAllPublisher(user))
			case "get_publishers":
				response = marshalPubSubResult(metrics.GetPublishers(user, parseTopic(argsByte)))
			case "get_subscribers":
				response = marshalPubSubResult(metrics.GetSubscribers(user, parseTopic(argsByte)))
			case "get_topic_statistics":
				response = marshalPubSubResult(metrics.GetTopicStatistics(user, parseTopic(argsByte)))
			case "get_inactive_topics":
				response = marshalPubSubResult(metrics.GetInactiveTopics(user))
			case "broadcast_to_topic":
				topic, message := parseTopicMessage(argsByte)
				response = marshalPubSubResult(metrics.BroadcastToTopic(user, topic, message))
			case "broadcast_everywhere":
				response = marshalPubSubResult(metrics.BroadcastEverywhere(user, parseMessage(argsByte)))
			case "history_of_topic":
				response = marshalPubSubResult(metrics.HistoryOfTopic(user, parseTopic(argsByte)))
			case "delete_topic":
				response = marshalPubSubResult(metrics.DeleteTopic(user, parseTopic(argsByte)))
			case "create_topic":
				response = marshalPubSubResult(metrics.CreateTopic(user, parseTopic(argsByte)))
			case "remove_user":
				portAddress, reason := parseRemoveUserArgs(argsByte)
				response = marshalPubSubResult(metrics.RemoveUser(user, portAddress, reason))
			case "get_total_keys":
				response = metrics.GetTotalKeys(stre)
			case "get_largest_keys":
				response = metrics.GetLargestKeys(stre, parseTopN(argsByte))
			case "get_smallest_keys":
				response = metrics.GetSmallestKeys(stre, parseTopN(argsByte))
			case "get_average_value_size":
				response = metrics.GetAverageValueSize(stre)
			case "get_namespaces":
				response = metrics.GetNamespaces(stre)
			case "get_ttl_metrics":
				response = metrics.GetTTLMetrics()
			case "get_expired_keys":
				response = metrics.GetExpiredKeys(stre)
			case "get_upcoming_expirations":
				response = metrics.GetUpcomingExpirations(stre, parseTopN(argsByte))
			case "get_keys_without_ttl":
				response = metrics.GetKeysWithoutTTL(stre)
			case "get_snapshot_count":
				response = metrics.GetSnapshotCount()
			case "get_latest_snapshot":
				response = metrics.GetLatestSnapshot()
			case "get_snapshot_size":
				response = metrics.GetSnapshotSize()
			case "get_persistence_health":
				response = metrics.GetPersistenceHealth()
			case "get_oldest_key":
				response = summarizeKeyMetadata(metrics.GetOldestKey(stre))
			case "get_newest_key":
				response = summarizeKeyMetadata(metrics.GetNewestKey(stre))
			case "get_most_accessed_key":
				response = summarizeKeyMetadata(metrics.GetMostAccessedKey(stre))
			case "get_least_accessed_key":
				response = summarizeKeyMetadata(metrics.GetLeastAccessedKey(stre))
			case "get_recently_modified_keys":
				response = summarizeKeyMetadataSlice(metrics.GetRecentlyModifiedKeys(stre, parseTopN(argsByte)))

			default:
				response = "No specialist response was produced."
			}

			if response == "" {
				response = "No specialist response was produced."
			}
			data := fmt.Sprintf("%v", response)
			fmt.Println("Data Appended: ", data)
			rawReplies = append(rawReplies, data)
			specialistReplies = append(specialistReplies, fmt.Sprintf("[%s]\n%s", tool.Function.Name, response))
		}
		session.AppendToolResults(rawReplies...)

		_, finalText, cerr := session.Run(ctx, client, agents.AGENT_USED, "Synthesize the worker outputs into the final user-facing answer. Interpret raw execution data, highlight the result that matters, and keep the reply concise unless the user requested deeper analysis.", nil, MasterInference)
		if cerr != nil {
			log.Fatalf("Master analysis request failed: %v", cerr)
		}

		if strings.TrimSpace(finalText) != "" {
			return finalText
		}

		return strings.Join(specialistReplies, "\n\n")
	}

	if assistantTextResponse != "" {
		return assistantTextResponse
	}

	return "\nroxai> I’m ready to help with the requested RoxKV task.\n"
}

func marshalMonitorResult(value any) string {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", value)
	}

	return string(data)
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
