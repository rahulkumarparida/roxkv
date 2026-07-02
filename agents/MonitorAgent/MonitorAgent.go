package monitoragent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/rahulkumarparida/roxkv/internal/metrics"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)

func MonitorAgent(query string, stre *store.MemoryAlloc, user *utils.NewClient, client *api.Client) {
	_ = stre

	ctx := context.Background()
	stream := false

	Messages := []api.Message{
		{
			Role: "system",
			Content: `You are an AI monitoring assistant.
		Use monitoring tools to collect live system information.
		Never estimate or invent metrics.
		After gathering data, explain it in a clear and concise way.`,
		},
		{
			Role:    "user",
			Content: query,
		},
	}

	req := &api.ChatRequest{
		Model:    "llama3.2:3b",
		Messages: Messages,
		Tools: []api.Tool{
			GetComputerUsageTool(),
			GetCPUUsageTool(),
			GetRAMUsageTool(),
			GetDiskUsageTool(),
			GetRuntimeStatsTool(),
			GetServerUptimeTool(),
			GetConnectedClientsTool(),
			GetCommandCountTool(),
		},
		Stream:  &stream,
		Options: MonitorInference,
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
		Messages = append(Messages, api.Message{
			Role:      "assistant",
			ToolCalls: toolCallsToExecute,
		})

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
				json.Unmarshal(argsByte, &args)

				path := strings.TrimSpace(args.Path)
				if path == "" {
					path = "/"
				}

				disk := metrics.GetDiskUsage(path)
				if disk.Err != nil {
					toolResult = fmt.Sprintf("failed to fetch disk usage for %q: %v", path, disk.Err)
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
				toolResult = "No Tools found"
			}

			Messages = append(Messages, api.Message{
				Role:    "tool",
				Content: toolResult,
			})
		}

		var finalResponse string

		secondReq := &api.ChatRequest{
			Model:    "llama3.2:3b",
			Messages: Messages,
			Stream:   &stream,
			Options:  MonitorInference,
		}

		secondErr := client.Chat(ctx, secondReq, func(resp api.ChatResponse) error {
			if resp.Message.Content != "" {
				finalResponse = resp.Message.Content
			}
			return nil
		})

		if secondErr != nil {
			log.Fatalf("Second Ollama API call failed: %v", secondErr)
		}

		if finalResponse != "" {
			user.Conn.Write([]byte("roxai> " + finalResponse + "\n"))
		}

		return
	}

	if assistantTextResponse != "" {
		user.Conn.Write([]byte("\nroxai> " + assistantTextResponse + "\n"))
	}
}

func marshalMonitorResult(value any) string {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", value)
	}

	return string(data)
}
