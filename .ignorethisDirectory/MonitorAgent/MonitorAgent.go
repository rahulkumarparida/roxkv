package monitoragent

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

const monitorSystemPrompt = `You are the RoxKV monitoring worker.
You are not a conversational assistant.
Use only the provided monitoring tools to execute the request.
Return raw tool results only.
Do not summarize, explain, recommend, or invent metrics.
If no available tool can satisfy the request, return a structured error.`

func MonitorAgent(query string, stre *store.MemoryAlloc, user *utils.NewClient, client *api.Client) string {
	_ = stre

	ctx := context.Background()
	session := agents.GetSession("monitoragent", monitorSystemPrompt, user)

	tools := []api.Tool{
		GetComputerUsageTool(),
		GetCPUUsageTool(),
		GetRAMUsageTool(),
		GetDiskUsageTool(),
		GetRuntimeStatsTool(),
		GetServerUptimeTool(),
		GetConnectedClientsTool(),
		GetCommandCountTool(),
	}

	toolCallsToExecute, assistantTextResponse, cerr := session.Run(ctx, client, agents.AGENT_USED, query, tools, MonitorInference)
	if cerr != nil {
		log.Fatalf("Monitor worker API call failed: %v", cerr)
	}

	if len(toolCallsToExecute) == 0 {
		if strings.TrimSpace(assistantTextResponse) == "" {
			return agents.StructuredError("no_suitable_tool", "The monitoring worker could not match the request to an available monitoring tool.")
		}
		return agents.StructuredError("no_suitable_tool", strings.TrimSpace(assistantTextResponse))
	}

	var rawResults []string
	var results []agents.ToolResult
	for _, tool := range toolCallsToExecute {
		argsByte, _ := json.Marshal(tool.Function.Arguments)
		var toolResult string

		switch tool.Function.Name {
		case "get_computer_usage":
			toolResult = marshalMonitorResult(metrics.GetComputerUsage())
		case "get_cpu_usage":
			toolResult = metrics.GetCPUUsage()
		case "get_ram_usage":
			toolResult = marshalMonitorResult(metrics.GetRAMUsage())
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
				toolResult = agents.StructuredError("disk_usage_failed", fmt.Sprintf("failed to fetch disk usage for %q: %v", path, disk.Err))
			} else {
				toolResult = marshalMonitorResult(map[string]any{
					"path":         path,
					"total_gb":     disk.Total,
					"free_gb":      disk.Free,
					"available_gb": disk.Avaliable,
				})
			}
		case "get_runtime_stats":
			toolResult = marshalMonitorResult(metrics.GetRuntimeStats())
		case "get_server_uptime":
			toolResult = metrics.GetServerUptime().String()
		case "get_connected_clients":
			toolResult = marshalMonitorResult(metrics.GetConnectedClients())
		case "get_command_count":
			toolResult = fmt.Sprintf("%d", metrics.GetCommandCount())
		default:
			toolResult = agents.StructuredError("unknown_tool", "The monitoring worker received an unsupported tool call.")
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

func marshalMonitorResult(value any) string {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", value)
	}

	return string(data)
}
