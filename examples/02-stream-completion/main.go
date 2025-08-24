package main

import (
	"context"
	"fmt"

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
			return nil // Continue streaming
			// This do not stop the streaming, it just returns the content
			//return &mu.ExitStreamCompletionError{Message: "❌ EXIT"} // This will stop the streaming
		})

	fmt.Println() // Ensure the output ends with a newline
	if err != nil {
		panic(err)
	}

}
