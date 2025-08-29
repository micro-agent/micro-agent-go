package mu

import (
	"encoding/json"
	"errors"

	"github.com/openai/openai-go/v2"
)

// 🚧 Work In Progress

// RunWithReasoning executes a chat completion with the provided messages.
// It sends the messages to the model and returns the first choice's content and reasoning.
//
// Parameters:
//   - Messages: The conversation messages to send to the model
//
// Returns:
//   - string: The content of the first choice from the model's response
//   - string: The reasoning content from the model's response
//   - error: Any error that occurred during the completion request or if no choices are returned
//
// This method temporarily sets the agent's Messages parameter and makes a synchronous
// completion request. It returns an error if the completion fails or if the response
// contains no choices.
func (agent *BasicAgent) RunWithReasoning(Messages []openai.ChatCompletionMessageParamUnion) (string, string, error) {
	// Preserve existing system messages from agent.Params
	existingSystemMessages := []openai.ChatCompletionMessageParamUnion{}
	for _, msg := range agent.Params.Messages {
		if msg.OfSystem != nil {
			existingSystemMessages = append(existingSystemMessages, msg)
		}
	}

	// Combine existing system messages with new messages
	agent.Params.Messages = append(existingSystemMessages, Messages...)
	completion, err := agent.Client.Chat.Completions.New(agent.ctx, agent.Params)

	if err != nil {
		return "", "", err
	}

	if len(completion.Choices) > 0 {
		jsonResponse := completion.Choices[0].Message.RawJSON()
		// extract the content of the reasoning_content field from the jsonResponse
		var reasoningContent struct {
			ReasoningContent string `json:"reasoning_content"`
		}
		err := json.Unmarshal([]byte(jsonResponse), &reasoningContent)
		if err != nil {
			return "", "", err
		}
		reasoning := reasoningContent.ReasoningContent
		// Trim whitespace from reasoning
		//reasoning = strings.TrimSpace(reasoning)

		content := completion.Choices[0].Message.Content

		// PHC - 2025-08-29
		// Append the full response as an assistant message to the agent's messages
		agent.Params.Messages = append(agent.Params.Messages, openai.AssistantMessage(content))

		return content, reasoning, nil
	} else {
		return "", "", errors.New("no choices found")
	}
}
