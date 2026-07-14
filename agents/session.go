package agents

import (
	"context"
	"fmt"
	"sync"

	"github.com/ollama/ollama/api"
	"github.com/rahulkumarparida/roxkv/internal/logger"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)

type AgentSession struct {
	messages []api.Message `json:"messages"`
}

var sessionRegistry sync.Map

func sessionKey(agentName string, user *utils.NewClient) string {
	return fmt.Sprintf("%s:%v", agentName, user.ID)
}

func GetSession(agentName, systemPrompt string, user *utils.NewClient) *AgentSession {
	key := sessionKey(agentName, user)
	if existing, ok := sessionRegistry.Load(key); ok {
		return existing.(*AgentSession)
	}

	session := &AgentSession{
		messages: []api.Message{
			{
				Role:    "system",
				Content: systemPrompt,
			},
		},
	}

	actual, _ := sessionRegistry.LoadOrStore(key, session)
	return actual.(*AgentSession)
}

func (s *AgentSession) Run(ctx context.Context, client *api.Client, model, query string, tools []api.Tool, options map[string]any) ([]api.ToolCall, string, error) {

	s.messages = append(s.messages, api.Message{
		Role:    "user",
		Content: query,
	})

	req := &api.ChatRequest{
		Model:    model,
		Messages: append([]api.Message(nil), s.messages...),
		Tools:    tools,
		Options:  options,
	}

	stream := false
	req.Stream = &stream

	var toolCalls []api.ToolCall
	var assistantText string

	err := client.Chat(ctx, req, func(resp api.ChatResponse) error {
		if len(resp.Message.ToolCalls) > 0 {
			toolCalls = resp.Message.ToolCalls
		}
		if resp.Message.Content != "" {
			assistantText = resp.Message.Content
		}
		return nil
	})
	if err != nil {
		return nil, "", err
	}

	if len(toolCalls) > 0 {
		s.messages = append(s.messages, api.Message{
			Role:      "assistant",
			ToolCalls: toolCalls,
		})
		return toolCalls, assistantText, nil
	}

	if assistantText != "" {
		s.messages = append(s.messages, api.Message{
			Role:    "assistant",
			Content: assistantText,
		})
	}
	logger.InfoLog("Assistant response: " + assistantText)

	return nil, assistantText, nil
}

func (s *AgentSession) AppendToolResults(results ...string) {

	for _, result := range results {
		s.messages = append(s.messages, api.Message{
			Role:    "tool",
			Content: result,
		})
	}
}
