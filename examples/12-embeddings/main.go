package main

import (
	"context"
	"encoding/json"

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
		mu.WithEmbeddingParams(
			openai.EmbeddingNewParams{
				Model: "ai/mxbai-embed-large",
			},
		),
	)
	if err != nil {
		panic(err)
	}

	embeddingVector, err := chatAgent.GenerateEmbeddingVector("I 💙 Docker Model Runner")

	if err != nil {
		panic(err)
	}
	jsonBytes, err := json.Marshal(embeddingVector)
	if err != nil {
		panic(err)
	}
	println("Embedding Vector:", string(jsonBytes))

}
