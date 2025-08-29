package main

import (
	"context"
	"fmt"

	"github.com/micro-agent/micro-agent-go/agent/msg"
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
			Model:       "ai/qwen2.5:1.5B-F16",
			Temperature: openai.Opt(0.0),
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("Your name is Bob. You are a helpful AI assistant."),
			},
		}),
	)
	if err != nil {
		panic(err)
	}

	for {
		content, _ := ui.SimplePrompt("🤖 (/bye to exit)>", "Type your command here...")

		if content.Input == "/bye" {
			ui.Println(ui.Green, "Goodbye!")
			break
		}

		_, err = chatAgent.RunStream(
			[]openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(content.Input),
			},
			func(content string) error {
				if content != "" {
					fmt.Print(content)
				}
				return nil // Continue streaming
			})

		fmt.Println() // Ensure the output ends with a newline
		if err != nil {
			panic(err)
		}

		ui.Println(ui.Red, "📝 Messages history / Conversational memory:")
		for i, message := range chatAgent.GetMessages() {
			printableMessage, err := msg.MessageToMap(message)
			if err != nil {
				ui.Printf(ui.Red, "Error converting message to map: %v\n", err)
				continue
			}
			ui.Print(ui.Cyan, "-", i, " ")
			ui.Print(ui.Orange, printableMessage["role"], ": ")
			ui.Println(ui.Cyan, printableMessage["content"])
		}
		fmt.Println()

	}
}
