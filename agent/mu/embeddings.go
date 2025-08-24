package mu

import "github.com/openai/openai-go/v2"

// GenerateEmbeddingVector creates a vector embedding for the given text content using the agent's embedding model
func (agent *Agent) GenerateEmbeddingVector(content string) ([]float64, error) {
	// Create embedding parameters using the agent's embedding parameters
	// params := openai.EmbeddingNewParams{
	// 	Model: agent.EmbeddingParams.Model,
	// 	Input: openai.EmbeddingNewParamsInputUnion{
	// 					OfString: openai.String(content),
	// 	},
	// }

	agent.EmbeddingParams.Input = openai.EmbeddingNewParamsInputUnion{
		OfString: openai.String(content),
	}
	// Use the client to create embeddings
	embeddingResponse, err := agent.Client.Embeddings.New(agent.ctx, agent.EmbeddingParams)
	if err != nil {
		return nil, err
	}

	return  embeddingResponse.Data[0].Embedding, nil
}
