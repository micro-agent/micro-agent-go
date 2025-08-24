// Package a2a provides experimental functionality for µ-agent.
//
// WARNING: This package is experimental and subject to change.
// The API may change or be removed in future versions without notice.
// Use at your own risk in production environments.
// NOTE: This is a partial implementation of the A2A protocol.
// IMPORTANT: This is a work in progress and may not cover all aspects of the A2A protocol.
package a2a

import (
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
