package kvagent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ollama/ollama/api"
	"github.com/rahulkumarparida/roxkv/agents"
	"github.com/rahulkumarparida/roxkv/internal/commands"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)

func KvAgent(query string, stre *store.MemoryAlloc, user *utils.NewClient, namespace *store.NameSpace, client *api.Client) {

	ctx := context.Background()
	stream := false

	Message := []api.Message{
		{
			Role: "system",
			Content: `You are an AI controller for an in-memory key-value database.
		Use tools whenever database access is required.
		Never fabricate database contents unless required.
		Answer normally if no tool is needed.`,
		},
		{
			Role:    "user",
			Content: query,
		},
	}
	// User will send a query
	req := &api.ChatRequest{
		Model:    agents.AGENT_USED,
		Messages: Message,
		Tools:    []api.Tool{GetKeyTool(), SetKeyTool(), KeysTool(), DeleteKeyTool(), SaveTool(), LoadTool()},
		Stream:   &stream,
		Options:  KvInference,
	}

	var toolCallsToExecute []api.ToolCall
	var assistantTextResponse string

	// LLM asks for tool execution
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
		// Executes the tools
		fmt.Printf("Model requested %d tool call(s).\n", len(toolCallsToExecute))
		Message = append(Message, api.Message{
			Role:      "assistant",
			ToolCalls: toolCallsToExecute,
		})

		var finalResponse []any 

		for _, tool := range toolCallsToExecute {

			argsByte, _ := json.Marshal(tool.Function.Arguments)

			switch tool.Function.Name {
			case "get":
				var args struct {
					Key string
				}
				json.Unmarshal(argsByte, &args)

				// Invoke native stock code
				result := store.GetKv(stre, strings.TrimSpace(args.Key))
				fmt.Printf("Executed get tool for: %s -> %v\n", args.Key, result)
				values := fmt.Sprintf("%v", result.Val)
				// Append tool result message
				// Message = append(Message, api.Message{
				// 	Role:    "tool",
				// 	Content: values,
				// })
				finalResponse = append(finalResponse, values)

			case "set":
				var args struct {
					Key   string
					Value string
				}
				json.Unmarshal(argsByte, &args)

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

				result := store.SetKv(stre, namespace, &item)

				values := fmt.Sprintf("%v", result)
				finalResponse = append(finalResponse, values)
				// Message = append(Message, api.Message{
				// 	Role:    "tool",
				// 	Content: values,
				// })

			case "del":
				var args struct {
					Key string
				}
				json.Unmarshal(argsByte, &args)

				result := store.DelKv(stre, strings.TrimSpace(args.Key))

				values := fmt.Sprintf("%v", result)

				// Message = append(Message, api.Message{
				// 	Role:    "tool",
				// 	Content: values,
				// })
				finalResponse = append(finalResponse, values)

			case "keys":
				result := store.KeyKv(stre)
				values := fmt.Sprintf("%v", result)
				finalResponse = append(finalResponse, values)

				// Message = append(Message, api.Message{
				// 	Role:    "tool",
				// 	Content: values,
				// })

			case "save":
				values := commands.SaveCommand(stre)
				finalResponse = append(finalResponse, values)

				// Message = append(Message, api.Message{
				// 	Role:    "tool",
				// 	Content: values,
				// })
			case "load":
				total := commands.LoaderCommand(stre, namespace)

				finalResponse = append(finalResponse, total)

				// Message = append(Message, api.Message{
				// 	Role:    "tool",
				// 	Content: strconv.Itoa(total),
				// })
			default:
				// Message = append(Message, api.Message{
				// 	Role:    "tool",
				// 	Content: "No Tools found",
				// })
				continue
			}

		}
		
		value := fmt.Sprintf("%v",finalResponse)

		user.Conn.Write([]byte("\nroxai> " + value + "\n"))

		return

	} else if assistantTextResponse != "" {
		// fmt.Fprintf(user.Conn, "%s\n", assistantTextResponse)
		user.Conn.Write([]byte("\nroxai> " + assistantTextResponse + "\n"))
		return
	}
}
