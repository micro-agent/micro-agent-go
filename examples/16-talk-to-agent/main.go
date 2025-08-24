package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/micro-agent/micro-agent-go/agent/experimental/a2a"
)

func main() {

	//ctx := context.Background()

	a2aClient := a2a.NewA2AClient("http://localhost:7777")
	agentCard, err := a2aClient.PingAgent()
	if err != nil {
		fmt.Printf("❌ Failed to ping agent: %v\n", err)
		return
	}
	fmt.Printf("🟢 Successfully pinged agent: %s - %s\n", agentCard.Name, agentCard.Description)

	jsonStrCrd, _ := a2a.AgentCardToJSONString(agentCard)
	fmt.Printf("🟢 AgentCard JSON:\n%s\n", jsonStrCrd)

	taskRequest := a2a.TaskRequest{
		ID:     uuid.NewString(),
		Method: "message/send",
		Params: a2a.AgentMessageParams{
			Message: a2a.AgentMessage{
				Role: "user",
				Parts: []a2a.TextPart{
					{
						Text: "What is the best pizza in the world?",
					},
				},
			},
			MetaData: map[string]any{
				"skill": "ask_for_something",
			},
		},
	}

	taskResponse, err := a2aClient.SendToAgent(taskRequest)
	if err != nil {
		fmt.Println("Error sending task request:", err)
		return
	}

	jsonTaskResponse, err := a2a.TaskResponseToJSONString(taskResponse)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("🟢 Task Response JSON:\n", jsonTaskResponse)

	fmt.Println("🟣 Task Response Text:", taskResponse.Result.History[0].Parts[0].Text)

}
