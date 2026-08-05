package agents

import (
	"context"
	"fmt"
	"sync"

	"github.com/rahulkumarparida/roxkv/agents/abstractor"
	"github.com/rahulkumarparida/roxkv/internal/logger"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)

type AgentSession struct {
	messages []abstractor.GenericMessage `json:"messages"`
}

var sessionRegistry sync.Map



func sessionKey(agentName string, user *utils.NewClient) string {
	return fmt.Sprintf("%s:%v", agentName, user.ID)
}
// 
func GetSession(agentName, systemPrompt string, user *utils.NewClient) *AgentSession {
	key := sessionKey(agentName, user)
	if existing, ok := sessionRegistry.Load(key); ok {
		return existing.(*AgentSession)
	}

	session := &AgentSession{
		messages: []abstractor.GenericMessage{
			{
				Role:    "system",
				Content: systemPrompt,
			},
		},
	}

	actual, _ := sessionRegistry.LoadOrStore(key, session)
	return actual.(*AgentSession)
}

func (s *AgentSession) Run(ctx context.Context, provider abstractor.Provider, model, query string, tools []abstractor.GenericToolDefinition, options map[string]any) ([]abstractor.GenericToolCall, string, error) {

	s.messages = append(s.messages, abstractor.GenericMessage{
		Role:    "user",
		Content: query,
	})

	// Get provider config (assuming we have one in context or global, but the provider handles it)
	// We'll pass a dummy config for now, or fetch the active config
	config := *abstractor.GetConfig()

	resp, err := provider.Chat(ctx, s.messages, tools, config)
	if err != nil {
		return nil, "", err
	}

	if len(resp.ToolCalls) > 0 {
		s.messages = append(s.messages, abstractor.GenericMessage{
			Role:      "assistant",
			ToolCalls: resp.ToolCalls,
		})
		return resp.ToolCalls, resp.Text, nil
	}

	if resp.Text != "" {
		s.messages = append(s.messages, abstractor.GenericMessage{
			Role:    "assistant",
			Content: resp.Text,
		})
	}
	logger.InfoLog("Assistant response: " + resp.Text)

	return nil, resp.Text, nil
}

func (s *AgentSession) AppendToolResults(results ...string) {

	for _, result := range results {
		s.messages = append(s.messages, abstractor.GenericMessage{
			Role:    "tool",
			Content: result,
		})
	}
}
