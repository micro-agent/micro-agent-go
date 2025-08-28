package mu

import (
	"encoding/json"
	"errors"

	"github.com/openai/openai-go/v2"
)

// 🚧 Work In Progress


// RunStreamWithReasoning. executes a streaming chat completion with the given messages.
// It streams the response content and reasoning in real-time by calling the provided callbacks for each chunk.
// The complete response content and reasoning are also accumulated and returned at the end.
//
// Parameters:
//   - Messages: The conversation messages to send to the model
//   - contentCallback: Function called for each content streaming chunk. Takes content string, returns error to stop streaming
//   - reasoningCallback: Function called for each reasoning streaming chunk. Takes reasoning string, returns error to stop streaming
//
// Returns:
//   - string: The complete accumulated response content from all chunks
//   - string: The complete accumulated reasoning content from all chunks
//   - error: Any error that occurred during streaming or from the callbacks
//
// The streaming stops early if:
//   - Either callback returns a non-nil error
//   - A stream error occurs
//   - Stream closing fails
func (agent *BasicAgent) RunStreamWithReasoning(Messages []openai.ChatCompletionMessageParamUnion, contentCallback func(content string) error, reasoningCallback func(reasoning string) error) (string, string, error) {
	// Preserve existing system messages from agent.Params
	existingSystemMessages := []openai.ChatCompletionMessageParamUnion{}
	for _, msg := range agent.Params.Messages {
		if msg.OfSystem != nil {
			existingSystemMessages = append(existingSystemMessages, msg)
		}
	}

	// Combine existing system messages with new messages
	agent.Params.Messages = append(existingSystemMessages, Messages...)
	stream := agent.Client.Chat.Completions.NewStreaming(agent.ctx, agent.Params)
	var response string
	var reasoning string
	var cbkRes error

	for stream.Next() {
		chunk := stream.Current()

		// Stream content chunk as it arrives
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			cbkRes = contentCallback(chunk.Choices[0].Delta.Content)
			response += chunk.Choices[0].Delta.Content
			if cbkRes != nil {
				var exitErr *ExitStreamCompletionError
				if errors.As(cbkRes, &exitErr) {
					break
				}
			}
		}

		// Extract and stream reasoning content if available
		if len(chunk.Choices) > 0 {
			jsonResponse := chunk.Choices[0].Delta.RawJSON()
			var reasoningContent struct {
				ReasoningContent string `json:"reasoning_content"`
			}
			err := json.Unmarshal([]byte(jsonResponse), &reasoningContent)
			if err == nil && reasoningContent.ReasoningContent != "" {
				//reasoningChunk := strings.TrimSpace(reasoningContent.ReasoningContent)
				reasoningChunk := reasoningContent.ReasoningContent

				if reasoningChunk != "" {
					cbkRes = reasoningCallback(reasoningChunk)
					reasoning += reasoningChunk
					if cbkRes != nil {
						var exitErr *ExitStreamCompletionError
						if errors.As(cbkRes, &exitErr) {
							break
						}
					}
				}
			}
		}
	}

	if cbkRes != nil {
		return response, reasoning, cbkRes
	}
	if err := stream.Err(); err != nil {
		return response, reasoning, err
	}
	if err := stream.Close(); err != nil {
		return response, reasoning, err
	}

	return response, reasoning, nil
}
