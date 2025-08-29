package mu

import (
	"context"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/shared"
)

// Agent is the interface for AI agents that can interact with OpenAI models and tools
type Agent interface {
	Run(Messages []openai.ChatCompletionMessageParamUnion) (string, error)
	RunStream(Messages []openai.ChatCompletionMessageParamUnion, callBack func(content string) error) (string, error)
	RunWithReasoning(Messages []openai.ChatCompletionMessageParamUnion) (string, string, error)
	RunStreamWithReasoning(Messages []openai.ChatCompletionMessageParamUnion, contentCallback func(content string) error, reasoningCallback func(reasoning string) error) (string, string, error)
	DetectToolCalls(messages []openai.ChatCompletionMessageParamUnion, toolCallBack func(functionName string, arguments string) (string, error)) (string, []string, string, error)
	DetectToolCallsStream(messages []openai.ChatCompletionMessageParamUnion, toolCallback func(functionName string, arguments string) (string, error), streamCallback func(content string) error) (string, []string, string, error)
	GenerateEmbeddingVector(content string) ([]float64, error)
	GetMessages() []openai.ChatCompletionMessageParamUnion
	SetMessages(messages []openai.ChatCompletionMessageParamUnion)
	GetResponseFormat() openai.ChatCompletionNewParamsResponseFormatUnion
	SetResponseFormat(format openai.ChatCompletionNewParamsResponseFormatUnion)
	GetName() string
	SetName(name string)
	GetModel() shared.ChatModel
	SetModel(model shared.ChatModel)
}

// BasicAgent represents a basic implementation of Agent with OpenAI client configuration and UI properties
type BasicAgent struct {
	ctx             context.Context
	Client          openai.Client
	Params          openai.ChatCompletionNewParams
	EmbeddingParams openai.EmbeddingNewParams
	Name            string
	Avatar          string
	Color           string // used for UI display
}

// AgentOption is a functional option for configuring BasicAgent instances
type AgentOption func(*BasicAgent)

// NewAgent creates a new Agent instance with the specified configuration.
// It uses the functional options pattern to configure the agent's client, parameters, and other settings.
//
// Parameters:
//   - ctx: Context for managing the agent's lifecycle and cancellation
//   - name: Human-readable name for the agent (used for identification/logging)
//   - options: Variable number of AgentOption functions to configure the agent
//
// Returns:
//   - Agent: Configured agent instance ready for use
//   - error: Always nil in current implementation, reserved for future validation
//
// Example usage:
//
//	agent, err := NewAgent(ctx, "ChatBot",
//	  WithClient(openaiClient),
//	  WithParams(completionParams),
//	)
func NewAgent(ctx context.Context, name string, options ...AgentOption) (Agent, error) {

	agent := &BasicAgent{}
	agent.Name = name
	agent.ctx = ctx
	// Apply all options
	for _, option := range options {
		option(agent)
	}

	return agent, nil
}

// WithClient is a functional option that sets the OpenAI client for an agent.
// The client handles the connection to the OpenAI API or compatible endpoints.
//
// Parameters:
//   - client: OpenAI client instance configured with API key, base URL, and other connection settings
//
// Returns:
//   - AgentOption: A function that applies the client to an Agent during construction
//
// Example usage:
//
//	client := openai.NewClient(option.WithAPIKey("your-api-key"))
//	agent := NewAgent(ctx, "MyAgent", WithClient(client))
func WithClient(client openai.Client) AgentOption {
	return func(a *BasicAgent) {
		a.Client = client
	}
}

// WithParams is a functional option that sets the chat completion parameters for an agent.
// This includes model settings, temperature, tools, messages, and other completion options.
//
// Parameters:
//   - params: OpenAI chat completion parameters including model, temperature, tools, messages, etc.
//
// Returns:
//   - AgentOption: A function that applies the parameters to an Agent during construction
//
// Example usage:
//
//	agent := NewAgent(ctx, "MyAgent", WithParams(openai.ChatCompletionNewParams{
//	  Model: "gpt-4",
//	  Temperature: openai.Opt(0.7),
//	  Tools: myTools,
//	}))
func WithParams(params openai.ChatCompletionNewParams) AgentOption {
	return func(a *BasicAgent) {
		a.Params = params
	}
}

// WithEmbeddingParams sets the embedding model parameters for the agent's vector generation
func WithEmbeddingParams(embeddingParams openai.EmbeddingNewParams) AgentOption {
	return func(a *BasicAgent) {
		a.EmbeddingParams = embeddingParams
	}
}

// GetMessages returns the messages from the agent's parameters
func (agent *BasicAgent) GetMessages() []openai.ChatCompletionMessageParamUnion {
	return agent.Params.Messages
}

// SetMessages sets the messages in the agent's parameters
func (agent *BasicAgent) SetMessages(messages []openai.ChatCompletionMessageParamUnion) {
	agent.Params.Messages = messages
}

// GetResponseFormat returns the response format from the agent's parameters
func (agent *BasicAgent) GetResponseFormat() openai.ChatCompletionNewParamsResponseFormatUnion {
	return agent.Params.ResponseFormat
}

// SetResponseFormat sets the response format in the agent's parameters
func (agent *BasicAgent) SetResponseFormat(format openai.ChatCompletionNewParamsResponseFormatUnion) {
	agent.Params.ResponseFormat = format
}

// GetName returns the name of the agent
func (agent *BasicAgent) GetName() string {
	return agent.Name
}

// SetName sets the name of the agent
func (agent *BasicAgent) SetName(name string) {
	agent.Name = name
}

// GetModel returns the model from the agent's parameters
func (agent *BasicAgent) GetModel() shared.ChatModel {
	return agent.Params.Model
}

// SetModel sets the model in the agent's parameters
func (agent *BasicAgent) SetModel(model shared.ChatModel) {
	agent.Params.Model = model
}
