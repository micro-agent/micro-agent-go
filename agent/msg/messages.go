package msg

import (
	"encoding/json"

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
