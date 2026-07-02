package kvagent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ollama/ollama/api"
	"github.com/rahulkumarparida/roxkv/internal/commands"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)


func KvAgent(query string,stre *store.MemoryAlloc,user *utils.NewClient,client *api.Client) {
	
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
				Role: "user",
				Content: query,
			},
		}
	// User will send a query
	req := &api.ChatRequest{
		Model: "llama3.2:3b",
		Messages: Message,
		Tools: []api.Tool{GetKeyTool(),SetKeyTool(),KeysTool(),DeleteKeyTool(),SaveTool(),LoadTool()},
		Stream: &stream,
		Options: KvInference,
	}

	var toolCallsToExecute []api.ToolCall
	var assistantTextResponse string

	// LLM asks for tool execution
	cerr := client.Chat(ctx,req,func(resp api.ChatResponse) error{
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
			Message = append(Message , api.Message{
				Role:      "assistant",
				ToolCalls: toolCallsToExecute,
			})
			

			for _, tool := range toolCallsToExecute {


				argsByte, _ := json.Marshal(tool.Function.Arguments)


				switch tool.Function.Name {
					case "get":
						var args struct {
							Key string
						}
						json.Unmarshal(argsByte, &args)

						// Invoke native stock code
						result := store.GetKv(stre,strings.TrimSpace(args.Key))
						fmt.Printf("Executed get tool for: %s -> %v\n", args.Key, result)
						values := fmt.Sprintf("%v", result.Val)
						// Append tool result message
						Message = append(Message, api.Message{
							Role:    "tool",
							Content: values,
						})


						
					case "set":
						var args struct {
							Key string
							Value string
						}
						json.Unmarshal(argsByte, &args)

						store.GetKv(stre,strings.TrimSpace(args.Key))

						item := store.Item{
							Key: strings.TrimSpace(args.Key),
							Val: strings.TrimSpace(args.Value),
							Ttl: time.Time{},
							UpdatedAt: time.Now(),
							LastAcessedBy: user,
							Size: len(args.Value),
						}

						result := store.SetKv(stre,&item,[]string{"",""})

						values := fmt.Sprintf("%v", result)

						Message = append(Message, api.Message{
							Role:    "tool",
							Content: values,
						})


					case "del":
						var args struct {
							Key string
						}
						json.Unmarshal(argsByte, &args)

						result := store.DelKv(stre,strings.TrimSpace(args.Key))

						values := fmt.Sprintf("%v", result)

						Message = append(Message, api.Message{
							Role:    "tool",
							Content: values,
						})


					case "keys":
						result := store.KeyKv(stre)
						values := fmt.Sprintf("%v", result)

						Message = append(Message, api.Message{
							Role:    "tool",
							Content: values,
						})

					case "save":
						values := commands.SaveCommand(stre)
						Message = append(Message, api.Message{
							Role:    "tool",
							Content: values,
						})
					case "load":
						total :=commands.LoaderCommand(stre)

						Message = append(Message, api.Message{
							Role:    "tool",
							Content: strconv.Itoa(total),
						})
					default:
						Message = append(Message, api.Message{
							Role:    "tool",
							Content: "No Tools found",
						})

				}
				
			}
			// Finnaly asks for a message response for the given query ans tools executed
			var finalResponse string

			secondReq := &api.ChatRequest{
				Model: "llama3.2:3b",
				Messages: Message,
				Stream: &stream,
				Options: KvInference,
			}

			secondErr := client.Chat(ctx, secondReq, func(resp api.ChatResponse) error {
				fmt.Printf("Executes after sending the final response %+v\n", Message)
				fmt.Printf("%+v\n", req.Messages)
				if resp.Message.Content != "" {
					finalResponse = resp.Message.Content
				}
				return nil
			})

			if secondErr != nil {
				log.Fatalf("Second Ollama API call failed: %v", secondErr)
			}

			// 2. Print the model's final conversational answer to the user
			if finalResponse != "" {
				user.Conn.Write([]byte("roxai> " + finalResponse + "\n"))
			}
			
			return

		}else if assistantTextResponse != "" {
		// fmt.Fprintf(user.Conn, "%s\n", assistantTextResponse)
			user.Conn.Write([]byte("\nroxai> "+assistantTextResponse+"\n"))
			return
		}
}





