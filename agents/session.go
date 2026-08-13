package agents

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/rahulkumarparida/roxkv/agents/abstractor"
	"github.com/rahulkumarparida/roxkv/internal/logger"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)

type AgentSession struct {
	mu       sync.Mutex
	messages []abstractor.GenericMessage `json:"messages"`
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
	s.mu.Lock()
	s.messages = append(s.messages, abstractor.GenericMessage{
		Role:    "user",
		Content: query,
	})
	msgSnapshot := make([]abstractor.GenericMessage, len(s.messages))
	copy(msgSnapshot, s.messages)
	s.mu.Unlock()

	logger.InfoLog("Executing via LLM failover: provider=" + fmt.Sprintf("%v", provider) + " model=" + model)
	resp, activeProv, err := abstractor.ExecuteWithFailover(ctx, provider, model, options, msgSnapshot, tools)
	if err != nil {
		logger.ErrorLog("LLM Manager Failover Error: " + err.Error())
		return nil, "", err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

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
	logger.InfoLog("Assistant response from " + activeProv + ": " + resp.Text)

	return nil, resp.Text, nil
}

func (s *AgentSession) AppendToolResults(results ...abstractor.GenericToolResult) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, res := range results {
		content := strings.Join(res.Content,",")
		s.messages = append(s.messages, abstractor.GenericMessage{
			Role:       "tool",
			Content:    content,
			ToolCallID: res.ID,
			ToolName:   res.Name,
		})
	}
}
