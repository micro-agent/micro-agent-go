package main

import (
	"context"

	"github.com/micro-agent/micro-agent-go/agent/mu"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

func main() {

	ctx := context.Background()
	// Initialize OpenAI client
	client := openai.NewClient(
		option.WithBaseURL("http://localhost:12434/engines/llama.cpp/v1"),
		option.WithAPIKey(""),
	)

	chatAgent, err := mu.NewAgent(ctx, "Bob",
		mu.WithClient(client),
		mu.WithParams(openai.ChatCompletionNewParams{
			Model:       "ai/qwen2.5:1.5B-F16",
			Temperature: openai.Opt(0.0),
		}),
	)
	if err != nil {
		panic(err)
	}

	response, err := chatAgent.Run([]openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage("Your name is Bob. You are a helpful AI assistant."),
		openai.UserMessage("Hello what is your name?"),
	})
	if err != nil {
		panic(err)
	}
	println("Response:", response)

}
