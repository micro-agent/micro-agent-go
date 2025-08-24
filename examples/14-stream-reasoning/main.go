package main

import (
	"context"
	"fmt"

	"github.com/micro-agent/micro-agent-go/agent/mu"
	"github.com/micro-agent/micro-agent-go/agent/ui"
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
			Model:       "hf.co/menlo/lucy-gguf:q8_0",
			Temperature: openai.Opt(0.0),
			Messages:    []openai.ChatCompletionMessageParamUnion{},
		}),
	)
	if err != nil {
		panic(err)
	}

	_, _, err = chatAgent.RunStreamWithReasoning(
		[]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("Your name is Bob. You are a helpful AI assistant."),
			openai.UserMessage("Who is Jean-Luc Picard?"),
		},
		func(content string) error {
			if content != "" {
				fmt.Print(content)
			}
			return nil // Continue streaming
		},
		func(reasoningContent string) error {
			if reasoningContent != "" {
				ui.Print(ui.Blue, reasoningContent)
			}
			return nil // Continue streaming
		},
	)

	fmt.Println() // Ensure the output ends with a newline
	if err != nil {
		panic(err)
	}

}
