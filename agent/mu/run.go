package mu

import (
	"errors"

	"github.com/openai/openai-go/v2"
)

// Run executes a chat completion with the provided messages.
// It sends the messages to the model and returns the first choice's content.
//
// Parameters:
//   - Messages: The conversation messages to send to the model
//
// Returns:
//   - string: The content of the first choice from the model's response
//   - error: Any error that occurred during the completion request or if no choices are returned
//
// This method temporarily sets the agent's Messages parameter and makes a synchronous
// completion request. It returns an error if the completion fails or if the response
// contains no choices.
func (agent *BasicAgent) Run(Messages []openai.ChatCompletionMessageParamUnion) (string, error) {
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
		return "", err
	}

	if len(completion.Choices) > 0 {
		return completion.Choices[0].Message.Content, nil
	} else {
		return "", errors.New("no choices found")

	}
}
