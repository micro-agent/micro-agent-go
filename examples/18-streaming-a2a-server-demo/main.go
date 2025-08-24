package main

import (
	"context"
	"fmt"

	"github.com/micro-agent/micro-agent-go/agent/experimental/a2a"
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

	agentCard := a2a.AgentCard{
		Name:        "Bob",
		Description: "A helpful assistant with streaming support and expertise in the Star Trek universe.",
		URL:         "http://localhost:8888",
		Version:     "1.0.0",
		//Capabilities: map[string]any{},
		Skills: []map[string]any{
			{
				"id":          "ask_for_something",
				"name":        "Ask for something",
				"description": "Bob is using a small language model to answer questions with streaming",
			},
			{
				"id":          "greetings",
				"name":        "Say greetings",
				"description": "Bob can say greetings to a person with emojis using streaming",
			},
		},
	}

	// Synchronous callback (for /task endpoint)
	agentCallBack := func(taskRequest a2a.TaskRequest) (a2a.TaskResponse, error) {

		fmt.Printf("🟢 Processing synchronous task request: %s\n", taskRequest.ID)
		// Extract user message
		userMessage := taskRequest.Params.Message.Parts[0].Text
		fmt.Printf("🔵 UserMessage: %s\n", userMessage)
		fmt.Printf("🟡 TaskRequest Metadata: %v\n", taskRequest.Params.MetaData)

		var systemMessage, userPrompt string

		switch taskRequest.Params.MetaData["skill"] {
		case "ask_for_something":
			systemMessage = "You are Bob, a simple A2A agent. You can answer questions."
			userPrompt = userMessage

		case "greetings":
			systemMessage = "You are Bob, a simple A2A agent. You can answer questions."
			userPrompt = "Greetings to " + userMessage + " with emojis and use his name."

		default:
			systemMessage = "You are Bob, a simple A2A agent. You can answer questions."
			userPrompt = "Be nice, and explain that " + fmt.Sprintf("%v", taskRequest.Params.MetaData["skill"]) + " is not a valid task ID."
		}

		answer, err := chatAgent.Run([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemMessage),
			openai.UserMessage(userPrompt),
		})
		if err != nil {
			fmt.Printf("❌ Error during chat completion: %v\n", err)
			return a2a.TaskResponse{}, err
		}

		fmt.Printf("🤖 Generated response: %s\n", answer)

		// Create response task
		responseTask := a2a.TaskResponse{
			ID:             taskRequest.ID,
			JSONRpcVersion: "2.0",
			Result: a2a.Result{
				Status: a2a.TaskStatus{
					State: "completed",
				},
				History: []a2a.AgentMessage{
					{
						Role: "assistant",
						Parts: []a2a.TextPart{
							{
								Text: answer,
								Type: "text",
							},
						},
					},
				},
				Kind:     "task",
				Metadata: map[string]any{},
			},
		}

		return responseTask, nil
	}

	// Streaming callback (for /stream endpoint)
	agentStreamCallback := func(taskRequest a2a.TaskRequest, streamFunc func(content string) error) error {

		fmt.Printf("🟢 Processing streaming task request: %s\n", taskRequest.ID)
		// Extract user message
		userMessage := taskRequest.Params.Message.Parts[0].Text
		fmt.Printf("🔵 UserMessage: %s\n", userMessage)
		fmt.Printf("🟡 TaskRequest Metadata: %v\n", taskRequest.Params.MetaData)

		var systemMessage, userPrompt string

		switch taskRequest.Params.MetaData["skill"] {
		case "ask_for_something":
			systemMessage = "You are Bob, a simple A2A agent. You can answer questions."
			userPrompt = userMessage

		case "greetings":
			systemMessage = "You are Bob, a simple A2A agent. You can answer questions."
			userPrompt = "Greetings to " + userMessage + " with emojis and use his name."

		default:
			systemMessage = "You are Bob, a simple A2A agent. You can answer questions."
			userPrompt = "Be nice, and explain that " + fmt.Sprintf("%v", taskRequest.Params.MetaData["skill"]) + " is not a valid task ID."
		}

		// Use RunStream instead of Run for streaming
		_, err := chatAgent.RunStream(
			[]openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(systemMessage),
				openai.UserMessage(userPrompt),
			},
			func(content string) error {
				if content != "" {
					fmt.Print(content) // Print to console for debugging
					return streamFunc(content) // Stream to client
				}
				return nil // Continue streaming
			})

		fmt.Println() // Ensure the output ends with a newline
		if err != nil {
			fmt.Printf("❌ Error during streaming chat completion: %v\n", err)
			return err
		}

		return nil
	}

	a2aServer := a2a.NewA2AServerWithStreaming(7777, agentCard, agentCallBack, agentStreamCallback)
	fmt.Println("🚀 Starting A2A server with streaming support on port 7777...")
	if err := a2aServer.Start(); err != nil {
		fmt.Printf("❌ Failed to start A2A server: %v\n", err)
	}

}