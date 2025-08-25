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

	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{
				"type": "string",
			},
			"capital": map[string]any{
				"type": "string",
			},
			"languages": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "string",
				},
			},
		},
		"required": []string{"name", "capital", "languages"},
	}

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "country_info",
		Description: openai.String("Notable information about a country in the world"),
		Schema:      schema,
		Strict:      openai.Bool(true),
	}

	chatAgent, err := mu.NewAgent(ctx, "Bob",
		mu.WithClient(client),
		mu.WithParams(openai.ChatCompletionNewParams{
			Model:       "ai/qwen2.5:1.5B-F16",
			Temperature: openai.Opt(0.0),
			ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
				OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
					JSONSchema: schemaParam,
				},
			},
		}),
	)
	if err != nil {
		panic(err)
	}

	response, err := chatAgent.Run([]openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(`
		Your name is Bob. 
		You are an assistant that answers questions about countries around the world.
		`),
		openai.UserMessage("Tell me about Canada."),
	})
	if err != nil {
		panic(err)
	}
	println("Response:", response)

}
