// Package a2a provides experimental functionality for µ-agent.
//
// WARNING: This package is experimental and subject to change.
// The API may change or be removed in future versions without notice.
// Use at your own risk in production environments.
// NOTE: This is a partial implementation of the A2A protocol.
// IMPORTANT: This is a work in progress and may not cover all aspects of the A2A protocol.
package a2a

import "encoding/json"

func TaskRequestToJSONString(taskRequest TaskRequest) (string, error) {
	jsonData, err := json.MarshalIndent(taskRequest, "", "    ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func AgentCardToJSONString(agentCard AgentCard) (string, error) {
	jsonData, err := json.MarshalIndent(agentCard, "", "    ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func TaskResponseToJSONString(taskResponse TaskResponse) (string, error) {
	jsonData, err := json.MarshalIndent(taskResponse, "", "    ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
