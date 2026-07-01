package structures

import (
	"encoding/json"
	"time"
)


type ThinkValue struct {
	// Value can be a bool or string
	Value interface{}
}

type ToolCallFunctionArguments struct {
	// contains filtered or unexported fields
					Type       string                 `json:"type"`
				Properties map[string]interface{} `json:"properties"`
				Required   []string               `json:"required"`

}

type ToolCallFunction struct {
	Name      string                    `json:"name"`
	Description string
	Arguments ToolCallFunctionArguments `json:"arguments"`
}

type ToolCall struct {
	ID       string           `json:"id"`
	Function ToolCallFunction `json:"function"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	// Thinking contains the text that was inside thinking tags in the
	// original model output when ChatRequest.Think is enabled.
	Thinking   string      `json:"thinking"`
	ToolCalls  []ToolCall  `json:"tool_calls"`
}

type ChatRequest struct {
	// Model is the model name, as in [GenerateRequest].
	Model string `json:"model"`

	// Messages is the messages of the chat - can be used to keep a chat memory.
	Messages []Message `json:"messages"`

	// Stream enables streaming of returned responses; true by default.
	Stream *bool `json:"stream"`

	// Format is the format to return the response in (e.g. "json").
	Format json.RawMessage `json:"format"`

	// Tools is an optional list of tools the model has access to.
	Tools []ToolCall `json:"tools"`

	// Think controls whether thinking/reasoning models will think before
	// responding. Can be a boolean (true/false) or a string ("high", "medium", "low")
	// for supported models.
	Think *ThinkValue `json:"think"`

	// Truncate is a boolean that, when set to true, truncates the chat history messages
	// if the rendered prompt exceeds the context length limit.
	Truncate *bool `json:"truncate"`


}


type ChatResponse struct {
	// Model is the model name that generated the response.
	Model string `json:"model"`

	// RemoteModel is the name of the upstream model that generated the response.
	RemoteModel string `json:"remote_model"`

	// RemoteHost is the URL of the upstream Ollama host that generated the response.
	RemoteHost string `json:"remote_host"`

	// CreatedAt is the timestamp of the response.
	CreatedAt time.Time `json:"created_at"`

	// Message contains the message or part of a message from the model.
	Message Message `json:"message"`

	// Done specifies if the response is complete.
	Done bool `json:"done"`

	// DoneReason is the reason the model stopped generating text.
	DoneReason string `json:"done_reason"`

}
