package main

import (
	"fmt"
	"time"

	"github.com/micro-agent/micro-agent-go/agent/experimental/a2a"
)

func main() {
	// Initialize the A2A client
	client := a2a.NewA2AClient("http://localhost:7777")

	// First, ping the agent to verify connection
	fmt.Println("🔍 Pinging agent...")
	agentCard, err := client.PingAgent()
	if err != nil {
		fmt.Printf("❌ Failed to ping agent: %v\n", err)
		return
	}

	fmt.Printf("✅ Connected to agent: %s\n", agentCard.Name)
	fmt.Printf("📝 Description: %s\n", agentCard.Description)
	fmt.Printf("🔧 Available skills: %v\n", len(agentCard.Skills))
	fmt.Println()

	// Create a task request
	taskRequest := a2a.TaskRequest{
		ID:             fmt.Sprintf("task-%d", time.Now().Unix()),
		JSONRpcVersion: "2.0",
		Method:         "message/send",
		Params: a2a.AgentMessageParams{
			Message: a2a.AgentMessage{
				Role: "user",
				Parts: []a2a.TextPart{
					{
						Text: "Tell me a story about a robot learning to cook pizza",
						Type: "text",
					},
				},
			},
			MetaData: map[string]any{
				"skill": "ask_for_something",
			},
		},
	}

	fmt.Printf("🚀 Sending streaming task request: %s\n", taskRequest.ID)
	fmt.Println("🌊 Streaming response:")

	// Stream callback to display content as it arrives
	streamCallback := func(content string) error {
		fmt.Print(content) // Print each chunk as it arrives
		return nil         // Continue streaming
	}

	// Send the streaming request
	response, err := client.SendToAgentStream(taskRequest, streamCallback)
	if err != nil {
		fmt.Printf("\n❌ Failed to send streaming request: %v\n", err)
		return
	}

	fmt.Println()
	fmt.Printf("✅ Task completed: %s\n", response.ID)
	fmt.Printf("🎯 Status: %s\n", response.Result.Status.State)
	
	if len(response.Result.History) > 0 {
		fullText := response.Result.History[0].Parts[0].Text
		fmt.Printf("📝 Total response length: %d characters\n", len(fullText))
	}

	fmt.Println("\n🏁 Streaming demo completed!")
}