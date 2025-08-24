# Micro Agent

**µAgent** is a Go tiny library that simplifies the creation of AI agents using the OpenAI API. See it as a lightweight wrapper around the [OpenAI Go SDK](https://github.com/openai/openai-go).

> AI Agent is only a pattern. **µAgent** only helps you implement it faster.

> You can use any LLM provider compatible with OpenAI API - I mainly work with [Docker Model Runner](https://docs.docker.com/ai/model-runner/).

## Quick Start Example

```go
func main() {

	ctx := context.Background()

    client := openai.NewClient(
		option.WithBaseURL("http://localhost:12434/engines/llama.cpp/v1"),
		option.WithAPIKey(""),
	)

	chatAgent, err := mu.NewAgent(ctx, "Bob",
		mu.WithClient(client),
		mu.WithParams(openai.ChatCompletionNewParams{
			Model:       "ai/qwen2.5:1.5B-F16",
			Temperature: openai.Opt(0.0),
			Messages:    []openai.ChatCompletionMessageParamUnion{},
		}),
	)
	if err != nil {
		panic(err)
	}

	_, err = chatAgent.RunStream(
		[]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("Your name is Bob. You are a helpful AI assistant."),
			openai.UserMessage("Hello what is your name? What can you do for me?"),
		},
		func(content string) error {
			if content != "" {
				fmt.Print(content)
			}
			return nil 
		})

	if err != nil {
		panic(err)
	}
}
```
