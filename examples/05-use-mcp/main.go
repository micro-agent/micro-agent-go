package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/micro-agent/micro-agent-go/agent/mu"
	"github.com/micro-agent/micro-agent-go/agent/tools"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

func main() {

	ctx := context.Background()

	client := openai.NewClient(
		option.WithBaseURL("http://localhost:12434/engines/llama.cpp/v1"),
		option.WithAPIKey(""),
	)

	mcpClient, err := tools.NewStreamableHttpMCPClient(ctx, "http://localhost:9011")
	if err != nil {
		panic(fmt.Errorf("failed to create MCP client: %v", err))
	}
	fmt.Println("MCP Client initialized successfully")
	toolsIndex := mcpClient.OpenAITools()
	for _, tool := range toolsIndex {
		fmt.Printf("Tool: %s - %s\n", tool.GetFunction().Name, tool.GetFunction().Description)
	}

	toolAgent, err := mu.NewAgent(ctx, "Bob",
		mu.WithClient(client),
		mu.WithParams(openai.ChatCompletionNewParams{
			Model:       "hf.co/menlo/jan-nano-gguf:q4_k_m",
			Temperature: openai.Opt(0.0),
			ToolChoice: openai.ChatCompletionToolChoiceOptionUnionParam{
				OfAuto: openai.String("auto"),
			},
			Tools:             toolsIndex,
			ParallelToolCalls: openai.Opt(false),
		}),
	)
	if err != nil {
		panic(err)
	}

	// Say "Exit" to stop the process
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage(`
			search a snippet for doing hello world,	
			then generate a markdow document to explain it
		`),
	}

	// Create executeFunction with MCP client option
	executeFn := createExecuteFunction(mcpClient)

	finishReason, results, assistantMessage, err := toolAgent.DetectToolCalls(messages, executeFn)
	if err != nil {
		panic(err)
	}
	fmt.Printf("✋ Finish Reason: %s\n", finishReason)
	fmt.Printf("📝 Results: %v\n", results)
	fmt.Printf("🤓 Assistant Message: %s\n", assistantMessage)

}

func createExecuteFunction(mcpClient *tools.MCPClient) func(string, string) (string, error) {
	return func(functionName string, arguments string) (string, error) {
		fmt.Printf("🟢 Executing function: %s with arguments: %s\n", functionName, arguments)

		// If MCP client is available, use it to execute the tool
		if mcpClient != nil {
			ctx := context.Background()
			result, err := mcpClient.CallTool(ctx, functionName, arguments)
			if err != nil {
				return "", fmt.Errorf("MCP tool execution failed: %v", err)
			}

			if len(result.Content) > 0 {
				// Take the first content item and return its text
				resultContent := result.Content[0].(mcp.TextContent).Text

				fmt.Printf("✅ Tool executed successfully, result: %s\n", resultContent)

				return fmt.Sprintf(`{"result": "%s"}`, resultContent), nil
			}
			return `{"result": "Tool executed successfully but returned no content"}`, nil
		}

		// Fallback to local execution if no MCP client
		var args map[string]interface{}
		if err := json.Unmarshal([]byte(arguments), &args); err != nil {
			return "", fmt.Errorf("failed to unmarshal arguments: %v", err)
		}
		fmt.Printf("Function arguments: %+v\n", args)
		return `{"result": "Function executed locally"}`, nil
	}
}
