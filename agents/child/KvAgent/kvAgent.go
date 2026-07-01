package child

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

	message := []api.Message{
		{
			Role: "user",
			Content: query,
		},
	}


	req := &api.ChatRequest{
		Model: "llama3.2:3b",
		Messages: message,
		Tools: []api.Tool{GetKeyTool(),SetKeyTool(),KeysTool(),DeleteKeyTool(),SaveTool(),LoadTool()},
		Stream: &stream,
	}

	var toolCallsToExecute []api.ToolCall
	var assistantTextResponse string

	cerr := client.Chat(ctx,req,func(resp api.ChatResponse) error{
			if len(resp.Message.ToolCalls) > 0 {
			toolCallsToExecute = resp.Message.ToolCalls
			}
			if resp.Message.Content != "" {
			assistantTextResponse += resp.Message.Content
			}
			return nil
		})


		if cerr != nil {
			log.Fatalf("First Ollama API call failed: %v", cerr)
		}

		if len(toolCallsToExecute) > 0 {
			
			fmt.Printf("Model requested %d tool call(s).\n", len(toolCallsToExecute))
			message = append(message , api.Message{
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
						message = append(message, api.Message{
							Role:    "tool",
							Content: values,
						})
						user.Conn.Write([]byte("roxai> "+values+"\n"))	
						
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
						user.Conn.Write([]byte("roxai> "+values+"\n"))

					case "del":
						var args struct {
							Key string
						}
						json.Unmarshal(argsByte, &args)

						result := store.DelKv(stre,strings.TrimSpace(args.Key))

						values := fmt.Sprintf("%v", result)

						user.Conn.Write([]byte("roxai> "+values+"\n"))

					case "keys":
						result := store.KeyKv(stre)
						values := fmt.Sprintf("%v", result)

						user.Conn.Write([]byte("roxai> "+values+"\n"))
					case "save":
						values := commands.SaveCommand(stre)
						user.Conn.Write([]byte("roxai> "+values+"\n"))
					case "load":
						total :=commands.LoaderCommand(stre)
						user.Conn.Write([]byte("roxai> "+ string(rune(total)) +"\n"))
					default:
						user.Conn.Write([]byte("roxai> "+"No tool found for the operration\n"))
						return 
						
				}
				
			}
				
			
		}
	
		if assistantTextResponse != "" {
		// fmt.Fprintf(user.Conn, "%s\n", assistantTextResponse)
			user.Conn.Write([]byte("\nroxai> "+assistantTextResponse+"\n"))
			return
		}
	


  

}





