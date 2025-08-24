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
			Model:       "hf.co/menlo/jan-nano-gguf:q4_k_m",
			Temperature: openai.Opt(0.0),
			Messages:    []openai.ChatCompletionMessageParamUnion{},
		}),
	)
	if err != nil {
		panic(err)
	}

	answer, err := chatAgent.RunStream(
		[]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("Your name is Bob. You are a helpful AI assistant."),
			openai.UserMessage("I want a hello world program in Go that prints 'Hello, World!' to the console."),
			openai.UserMessage("I want only the source code, nothing else. No explanations."),

			//openai.UserMessage("Who is Jean-Luc Picard?"),

		},
		func(content string) error {
			if content != "" {
				fmt.Print(content)
			}
			return nil // Continue streaming
		})

	if err != nil {
		fmt.Printf("Error during streaming: %v\n", err)
	}

	fmt.Println() // Ensure the output ends with a newline

	ui.PrintMarkdown(answer)

	// Copy answer to clipboard
	err = ui.CopyToClipboard(answer)
	if err != nil {
		fmt.Printf("Error copying to clipboard: %v\n", err)
	} else {
		fmt.Println("Answer copied to clipboard!")
	}

}
