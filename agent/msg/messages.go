package msg

import (
	"encoding/json"
	"fmt"

	"github.com/micro-agent/micro-agent-go/agent/mu"
	"github.com/micro-agent/micro-agent-go/agent/ui"
	"github.com/openai/openai-go/v2"
)

// MessageToMap converts an OpenAI chat message to a map with string keys and values
func MessageToMap(message openai.ChatCompletionMessageParamUnion) (map[string]string, error) {
	jsonData, err := message.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, err
	}

	stringMap := make(map[string]string)
	for key, value := range result {
		if str, ok := value.(string); ok {
			stringMap[key] = str
		}
	}

	return stringMap, nil
}

// MessagesToSlice converts a slice of OpenAI chat messages to a slice of string maps
func MessagesToSlice(messages []openai.ChatCompletionMessageParamUnion) ([]map[string]string, error) {
	result := make([]map[string]string, len(messages))

	for i, message := range messages {
		messageMap, err := MessageToMap(message)
		if err != nil {
			return nil, err
		}
		result[i] = messageMap
	}

	return result, nil
}

func DisplayHistory(selectedAgent mu.Agent) {
	// remove the /debug part from the input
	fmt.Println()
	ui.Println(ui.Red, "📝 Messages history / Conversational memory:")
	for i, message := range selectedAgent.GetMessages() {
		printableMessage, err := MessageToMap(message)
		if err != nil {
			ui.Printf(ui.Red, "Error converting message to map: %v\n", err)
			continue
		}
		ui.Print(ui.Cyan, "-", i, " ")
		ui.Print(ui.Orange, printableMessage["role"], ": ")
		ui.Println(ui.Cyan, printableMessage["content"])
	}
	ui.Println(ui.Red, "📝 End of the messages")
	fmt.Println()
}
