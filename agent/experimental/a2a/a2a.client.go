// Package a2a provides experimental functionality for µ-agent.
//
// WARNING: This package is experimental and subject to change.
// The API may change or be removed in future versions without notice.
// Use at your own risk in production environments.
// NOTE: This is a partial implementation of the A2A protocol.
// IMPORTANT: This is a work in progress and may not cover all aspects of the A2A protocol.
package a2a

import (
	"bufio"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type A2AClient struct {
	agentBaseURL string
}

func NewA2AClient(agentBaseURL string) *A2AClient {
	return &A2AClient{
		agentBaseURL: strings.TrimRight(agentBaseURL, "/"),
	}
}

func (a2acli *A2AClient) PingAgent() (AgentCard, error) {
	resp, err := http.Get(a2acli.agentBaseURL + "/.well-known/agent.json")
	if err != nil {
		return AgentCard{}, err
	}
	defer resp.Body.Close()

	var agentCard AgentCard
	if resp.StatusCode == http.StatusOK {
		if err := json.NewDecoder(resp.Body).Decode(&agentCard); err != nil {
			return AgentCard{}, err
		}
		return agentCard, nil
	} else {
		return agentCard, errors.New("failed to ping agent: " + resp.Status)
	}
}

func (a2acli *A2AClient) SendToAgent(taskRequest TaskRequest) (TaskResponse, error) {
	jsonTaskRequest, err := TaskRequestToJSONString(taskRequest)
	if err != nil {
		return TaskResponse{}, err
	}

	resp, err := http.Post(a2acli.agentBaseURL+"/", "application/json", strings.NewReader(jsonTaskRequest))
	if err != nil {
		return TaskResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TaskResponse{}, errors.New("failed to send task request: " + resp.Status)
	}

	var taskResponse TaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&taskResponse); err != nil {
		return TaskResponse{}, err
	}

	return taskResponse, nil
}

// SendToAgentStream sends a task request to the agent and streams the response
// streamCallback is called for each chunk of content received
// Returns the complete response and any error
func (a2acli *A2AClient) SendToAgentStream(taskRequest TaskRequest, streamCallback func(content string) error) (TaskResponse, error) {
	jsonTaskRequest, err := TaskRequestToJSONString(taskRequest)
	if err != nil {
		return TaskResponse{}, err
	}

	req, err := http.NewRequest("POST", a2acli.agentBaseURL+"/stream", strings.NewReader(jsonTaskRequest))
	if err != nil {
		return TaskResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return TaskResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TaskResponse{}, errors.New("failed to send streaming task request: " + resp.Status)
	}

	var finalResponse TaskResponse
	var fullContent strings.Builder

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		
		// Parse Server-Sent Events
		if strings.HasPrefix(line, "data: ") {
			jsonData := strings.TrimPrefix(line, "data: ")
			
			// Skip empty data lines
			if strings.TrimSpace(jsonData) == "" {
				continue
			}

			// Try to parse as chunk response
			var chunkResponse map[string]interface{}
			if err := json.Unmarshal([]byte(jsonData), &chunkResponse); err != nil {
				continue // Skip malformed JSON
			}

			// Handle streaming chunks
			if chunkType, exists := chunkResponse["type"]; exists && chunkType == "chunk" {
				if content, exists := chunkResponse["content"]; exists {
					if contentStr, ok := content.(string); ok {
						fullContent.WriteString(contentStr)
						if streamCallback != nil {
							if err := streamCallback(contentStr); err != nil {
								return TaskResponse{}, err
							}
						}
					}
				}
			} else {
				// Try to parse as final TaskResponse
				var taskResponse TaskResponse
				if err := json.Unmarshal([]byte(jsonData), &taskResponse); err == nil {
					finalResponse = taskResponse
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return TaskResponse{}, err
	}

	// If we don't have a final response, create one with the accumulated content
	if finalResponse.ID == "" {
		finalResponse = TaskResponse{
			ID:             taskRequest.ID,
			JSONRpcVersion: "2.0",
			Result: Result{
				Status: TaskStatus{
					State: "completed",
				},
				History: []AgentMessage{
					{
						Role: "assistant",
						Parts: []TextPart{
							{
								Text: fullContent.String(),
								Type: "text",
							},
						},
					},
				},
				Kind:     "task",
				Metadata: map[string]any{},
			},
		}
	}

	return finalResponse, nil
}
