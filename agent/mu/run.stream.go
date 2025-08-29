package mu

import (
	"errors"

	"github.com/openai/openai-go/v2"
)

// RunStream executes a streaming chat completion with the given messages.
// It streams the response content in real-time by calling the provided callback for each chunk.
// The complete response is also accumulated and returned at the end.
//
// Parameters:
//   - Messages: The conversation messages to send to the model
//   - callBack: Function called for each streaming chunk. Takes content string, returns error to stop streaming
//
// Returns:
//   - string: The complete accumulated response content from all chunks
//   - error: Any error that occurred during streaming or from the callback
//
// The streaming stops early if:
//   - The callback returns a non-nil error
//   - A stream error occurs
//   - Stream closing fails
func (agent *BasicAgent) RunStream(Messages []openai.ChatCompletionMessageParamUnion, callBack func(content string) error) (string, error) {
	// Preserve existing system messages from agent.Params
	// existingSystemMessages := []openai.ChatCompletionMessageParamUnion{}
	// for _, msg := range agent.Params.Messages {
	// 	if msg.OfSystem != nil {
	// 		existingSystemMessages = append(existingSystemMessages, msg)
	// 	}
	// }

	// Combine existing system messages with new messages
	agent.Params.Messages = append(agent.Params.Messages, Messages...)
	stream := agent.Client.Chat.Completions.NewStreaming(agent.ctx, agent.Params)
	var response string
	var cbkRes error

	for stream.Next() {
		chunk := stream.Current()
		// Stream each chunk as it arrives
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			cbkRes = callBack(chunk.Choices[0].Delta.Content)
			response += chunk.Choices[0].Delta.Content
		}

		// if cbkRes != nil {
		// 	break
		// }

		if cbkRes != nil {
			var exitErr *ExitStreamCompletionError
			if errors.As(cbkRes, &exitErr) {
				break
			}
		}

	}
	if cbkRes != nil {
		return response, cbkRes
	}
	if err := stream.Err(); err != nil {
		return response, err
	}
	if err := stream.Close(); err != nil {
		return response, err
	}

	// PHC - 2025-08-29
	// Append the full response as an assistant message to the agent's messages
	agent.Params.Messages = append(agent.Params.Messages, openai.AssistantMessage(response))

	return response, nil
}
