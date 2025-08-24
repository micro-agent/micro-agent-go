package mu

import (
	"context"
	"github.com/openai/openai-go/v2"
)

// Agent represents an AI agent with OpenAI client configuration and UI properties
type Agent struct {
	ctx             context.Context
	Client          openai.Client
	Params          openai.ChatCompletionNewParams
	EmbeddingParams openai.EmbeddingNewParams
	Name            string
	Avatar          string
	Color           string // used for UI display
}

// AgentOption is a functional option for configuring Agent instances
type AgentOption func(*Agent)

// NewAgent creates a new Agent instance with the specified configuration.
// It uses the functional options pattern to configure the agent's client, parameters, and other settings.
//
// Parameters:
//   - ctx: Context for managing the agent's lifecycle and cancellation
//   - name: Human-readable name for the agent (used for identification/logging)
//   - options: Variable number of AgentOption functions to configure the agent
//
// Returns:
//   - *Agent: Configured agent instance ready for use
//   - error: Always nil in current implementation, reserved for future validation
//
// Example usage:
//
//	agent, err := NewAgent(ctx, "ChatBot",
//	  WithClient(openaiClient),
//	  WithParams(completionParams),
//	)
func NewAgent(ctx context.Context, name string, options ...AgentOption) (*Agent, error) {

	agent := &Agent{}
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
	return func(a *Agent) {
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
	return func(a *Agent) {
		a.Params = params
	}
}

// WithEmbeddingParams sets the embedding model parameters for the agent's vector generation
func WithEmbeddingParams(embeddingParams openai.EmbeddingNewParams) AgentOption {
	return func(a *Agent) {
		a.EmbeddingParams = embeddingParams
	}
}
