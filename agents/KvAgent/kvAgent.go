package kvagent

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/ollama/ollama/api"
	"github.com/rahulkumarparida/roxkv/agents"
	"github.com/rahulkumarparida/roxkv/internal/commands"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)

const kvSystemPrompt = `You are the RoxKV key-value worker.
You are not a conversational assistant.
Use only the provided key-value tools to execute requests against the in-memory store.
Return raw tool results only.
Do not explain, summarize, speculate, or invent database contents.
If no available tool can satisfy the request, return a structured error.`

func KvAgent(query string, stre *store.MemoryAlloc, user *utils.NewClient, namespace *store.NameSpace, client *api.Client) string {
	ctx := context.Background()
	session := agents.GetSession("kvagent", kvSystemPrompt, user)
	tools := []api.Tool{GetKeyTool(), SetKeyTool(), KeysTool(), DeleteKeyTool(), SaveTool(), LoadTool()}

	toolCallsToExecute, assistantTextResponse, cerr := session.Run(ctx, client, agents.AGENT_USED, query, tools, KvInference)
	if cerr != nil {
		log.Fatalf("KV worker API call failed: %v", cerr)
	}

	if len(toolCallsToExecute) == 0 {
		if strings.TrimSpace(assistantTextResponse) == "" {
			return agents.StructuredError("no_suitable_tool", "The key-value worker could not match the request to an available key-value tool.")
		}
		return agents.StructuredError("no_suitable_tool", strings.TrimSpace(assistantTextResponse))
	}

	var rawResults []string
	var results []agents.ToolResult
	for _, tool := range toolCallsToExecute {
		argsByte, _ := json.Marshal(tool.Function.Arguments)
		var toolResult any

		switch tool.Function.Name {
		case "get":
			var args struct {
				Key string `json:"key"`
			}
			_ = json.Unmarshal(argsByte, &args)
			result := store.GetKv(stre, strings.TrimSpace(args.Key))
			toolResult = result.Val
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
			toolResult = store.SetKv(stre, namespace, &item)
		case "del":
			var args struct {
				Key string `json:"key"`
			}
			_ = json.Unmarshal(argsByte, &args)
			toolResult = store.DelKv(stre, strings.TrimSpace(args.Key))
		case "keys":
			toolResult = store.KeyKv(stre)
		case "save":
			toolResult = commands.SaveCommand(stre)
		case "load":
			toolResult = commands.LoaderCommand(stre, namespace)
		default:
			toolResult = map[string]any{
				"status":  "error",
				"error":   "unknown_tool",
				"message": "The key-value worker received an unsupported tool call.",
			}
		}

		rawBytes, _ := json.Marshal(toolResult)
		rawResults = append(rawResults, string(rawBytes))
		results = append(results, agents.ToolResult{
			Tool:   tool.Function.Name,
			Output: toolResult,
		})
	}

	session.AppendToolResults(rawResults...)
	return agents.StructuredSuccess(results)
}
