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
			Temperature: openai.Opt(0.5),
		}),
	)
	if err != nil {
		panic(err)
	}

	answer, reasoning, err := chatAgent.RunWithReasoning([]openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage("Your name is Bob. You are a helpful AI assistant."),
		openai.UserMessage("Who is Jean-Luc Picard?"),
	})
	if err != nil {
		panic(err)
	}

	ui.Println(ui.Blue, "🧠 Reasoning:\n", reasoning)
	fmt.Println()
	ui.Println(ui.Green, "🤖 Response:\n", answer)

}
