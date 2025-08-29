package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/micro-agent/micro-agent-go/agent/helpers"
	"github.com/micro-agent/micro-agent-go/agent/mu"
	"github.com/micro-agent/micro-agent-go/agent/rag"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

var store rag.MemoryVectorStore

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

	// -------------------------------------------------
	// Create a vector store
	// -------------------------------------------------
	store = rag.MemoryVectorStore{
		Records: make(map[string]rag.VectorRecord),
	}

	contextInstructionsContent, err := helpers.ReadTextFile("./sorcerer_background_and_personality.md")
	chunks := rag.SplitMarkdownBySections(contextInstructionsContent)
	for idx, chunk := range chunks {
		fmt.Println("🔶 Chunk", idx, ":", chunk)
		embeddingVector, err := chatAgent.GenerateEmbeddingVector(chunk)
		if err != nil {
			panic(err)
		}
		_, errSave := store.Save(rag.VectorRecord{
			Prompt:    chunk,
			Embedding: embeddingVector,
		})

		if errSave != nil {
			fmt.Println("😡:", errSave)
		}
		fmt.Println("✅ Chunk", idx, "saved with embedding:", len(embeddingVector))

	}

	fmt.Println("📝 Total records in the vector store:", len(store.Records))
	//fmt.Println(store.Records)

	fmt.Println(strings.Repeat("-", 80))
	question := "What is your name?"
	fmt.Printf("🔍 Searching for similar chunks to '%s'\n", question)
	questionEmbeddingVector, err := chatAgent.GenerateEmbeddingVector(question)

	questionRecord := rag.VectorRecord{
		Embedding: questionEmbeddingVector,
	}

	similarities, _ := store.SearchTopNSimilarities(questionRecord, 0.5, 2)

	fmt.Println("📝 Similarities found:", len(similarities))

	for _, similarity := range similarities {
		fmt.Println("✅ CosineSimilarity:", similarity.CosineSimilarity, "Chunk:", similarity.Prompt)
	}

	fmt.Println(strings.Repeat("-", 80))
	question = "Tell me about your family"
	fmt.Printf("🔍 Searching for similar chunks to '%s'\n", question)
	questionEmbeddingVector, err = chatAgent.GenerateEmbeddingVector(question)

	questionRecord = rag.VectorRecord{
		Embedding: questionEmbeddingVector,
	}

	similarities, _ = store.SearchTopNSimilarities(questionRecord, 0.5, 2)

	fmt.Println("📝 Similarities found:", len(similarities))

	for _, similarity := range similarities {
		fmt.Println("✅ CosineSimilarity:", similarity.CosineSimilarity, "Chunk:", similarity.Prompt)
	}
}
