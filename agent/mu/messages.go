package mu

import (
	"encoding/json"
	"github.com/openai/openai-go/v2"
)

// GetMessages returns the messages from the agent's parameters
func (agent *BasicAgent) GetMessages() []openai.ChatCompletionMessageParamUnion {
	return agent.Params.Messages
}

// GetFirstNMessages returns the first n messages from the agent's message list
func (agent *BasicAgent) GetFirstNMessages(n int) []openai.ChatCompletionMessageParamUnion {
	if n <= 0 {
		return []openai.ChatCompletionMessageParamUnion{}
	}
	messagesLen := len(agent.Params.Messages)
	if n >= messagesLen {
		return agent.Params.Messages
	}
	return agent.Params.Messages[:n]
}

// GetLastNMessages returns the last n messages from the agent's message list
func (agent *BasicAgent) GetLastNMessages(n int) []openai.ChatCompletionMessageParamUnion {
	if n <= 0 {
		return []openai.ChatCompletionMessageParamUnion{}
	}
	messagesLen := len(agent.Params.Messages)
	if n >= messagesLen {
		return agent.Params.Messages
	}
	return agent.Params.Messages[messagesLen-n:]
}

// GetFirstMessage returns the first message from the agent's message list
// Returns the message and true if found, or nil and false if the list is empty
func (agent *BasicAgent) GetFirstMessage() (openai.ChatCompletionMessageParamUnion, bool) {
	if len(agent.Params.Messages) == 0 {
		var zero openai.ChatCompletionMessageParamUnion
		return zero, false
	}
	return agent.Params.Messages[0], true
}

// GetLastMessage returns the last message from the agent's message list
// Returns the message and true if found, or nil and false if the list is empty
func (agent *BasicAgent) GetLastMessage() (openai.ChatCompletionMessageParamUnion, bool) {
	if len(agent.Params.Messages) == 0 {
		var zero openai.ChatCompletionMessageParamUnion
		return zero, false
	}
	return agent.Params.Messages[len(agent.Params.Messages)-1], true
}

// GetMessageByIndex returns the message at the specified index
// Returns the message and true if found, or nil and false if the index is out of bounds
func (agent *BasicAgent) GetMessageByIndex(index int) (openai.ChatCompletionMessageParamUnion, bool) {
	if index < 0 || index >= len(agent.Params.Messages) {
		var zero openai.ChatCompletionMessageParamUnion
		return zero, false
	}
	return agent.Params.Messages[index], true
}

// ReplaceMessageByIndex replaces the message at the specified index
// Returns true if the replacement was successful, false if the index is out of bounds
func (agent *BasicAgent) ReplaceMessageByIndex(index int, message openai.ChatCompletionMessageParamUnion) bool {
	if index < 0 || index >= len(agent.Params.Messages) {
		return false
	}
	agent.Params.Messages[index] = message
	return true
}

// ToJSON returns a JSON representation of the messages list
// Returns the JSON string and an error if marshaling fails
func (agent *BasicAgent) ToJSON() (string, error) {
	jsonBytes, err := json.Marshal(agent.Params.Messages)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// ToPrettyJSON returns a prettified JSON representation of the messages list
// Returns the formatted JSON string and an error if marshaling fails
func (agent *BasicAgent) ToPrettyJSON() (string, error) {
	jsonBytes, err := json.MarshalIndent(agent.Params.Messages, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// SetMessages sets the messages in the agent's parameters
func (agent *BasicAgent) SetMessages(messages []openai.ChatCompletionMessageParamUnion) {
	agent.Params.Messages = messages
}

// AddMessage adds a single message to the agent's message list
func (agent *BasicAgent) AddMessage(message openai.ChatCompletionMessageParamUnion) {
	agent.Params.Messages = append(agent.Params.Messages, message)
}

// AddMessages adds multiple messages to the agent's message list
func (agent *BasicAgent) AddMessages(messages []openai.ChatCompletionMessageParamUnion) {
	agent.Params.Messages = append(agent.Params.Messages, messages...)
}

// PrependMessage adds a message at the beginning of the agent's message list
func (agent *BasicAgent) PrependMessage(message openai.ChatCompletionMessageParamUnion) {
	agent.Params.Messages = append([]openai.ChatCompletionMessageParamUnion{message}, agent.Params.Messages...)
}

// PrependMessages adds multiple messages at the beginning of the agent's message list
func (agent *BasicAgent) PrependMessages(messages []openai.ChatCompletionMessageParamUnion) {
	agent.Params.Messages = append(messages, agent.Params.Messages...)
}

// ResetMessages clears all messages in the agent's parameters
func (agent *BasicAgent) ResetMessages() {
	agent.Params.Messages = nil
}

// RemoveLastMessage removes the last message from the agent's message list
func (agent *BasicAgent) RemoveLastMessage() {
	if len(agent.Params.Messages) > 0 {
		agent.Params.Messages = agent.Params.Messages[:len(agent.Params.Messages)-1]
	}
}

// RemoveLastNMessages removes the last n messages from the agent's message list
func (agent *BasicAgent) RemoveLastNMessages(n int) {
	if n <= 0 {
		return
	}
	messagesLen := len(agent.Params.Messages)
	if n >= messagesLen {
		agent.Params.Messages = nil
	} else {
		agent.Params.Messages = agent.Params.Messages[:messagesLen-n]
	}
}

// RemoveFirstMessage removes the first message from the agent's message list
func (agent *BasicAgent) RemoveFirstMessage() {
	if len(agent.Params.Messages) > 0 {
		agent.Params.Messages = agent.Params.Messages[1:]
	}
}
