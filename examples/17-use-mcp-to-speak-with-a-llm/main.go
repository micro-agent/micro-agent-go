package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/micro-agent/micro-agent-go/agent/mu"
	"github.com/micro-agent/micro-agent-go/agent/tools"
	"github.com/micro-agent/micro-agent-go/agent/ui"

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

	ui.Println(ui.Purple, "MCP Client initialized successfully")
	toolsIndex := mcpClient.OpenAIToolsWithFilter([]string{"ask_agent"})
	for _, tool := range toolsIndex {
		ui.Printf(ui.Magenta, "Tool: %s - %s\n", tool.GetFunction().Name, tool.GetFunction().Description)
	}

	toolAgent, err := mu.NewAgent(ctx, "Sam",
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

	content, _ := ui.SimplePrompt("🤖 (/bye to exit)>", "Type your command here...")

	// Say "Exit" to stop the process
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(`
			Your name is Sam.
			You are a helpful assistant.
		`),
		openai.UserMessage(content.Input),
	}

	// Stream callback for real-time content display
	streamCallback := func(thinkingCtrl, streamingCtrl *ui.ThinkingController) func(string) error {

		return func(content string) error {
			if thinkingCtrl.IsStarted() {
				thinkingCtrl.Stop()
				streamingCtrl.Start(ui.Green, "Streaming...")
			}

			if content != "" {
				cleanContent := strings.ReplaceAll("⏳ "+content, "\n", "")
				// Pad or truncate cleanContent to exactly 20 characters
				if len(cleanContent) > 20 {
					cleanContent = cleanContent[:20]
				} else {
					cleanContent = cleanContent + strings.Repeat(".", 20-len(cleanContent))
				}
				streamingCtrl.UpdateMessage("Stream Completion..." + cleanContent)
			}
			return nil
		}
	}

	thinkingCtrl := ui.NewThinkingController()
	thinkingCtrl.Start(ui.Blue, "Tools detection.....")
	streamingCtrl := ui.NewThinkingController()

	// Create executeFunction with MCP client option
	// Tool execution callback
	executeFn := executeFunction(mcpClient, thinkingCtrl)

	_, _, assistantMessage, err := toolAgent.DetectToolCallsStream(messages, executeFn, streamCallback(thinkingCtrl, streamingCtrl))
	if err != nil {
		panic(err)
	}

	thinkingCtrl.Stop()

	fmt.Println()
	fmt.Println()

	ui.PrintMarkdown(assistantMessage)

}

func executeFunction(mcpClient *tools.MCPClient, thinkingCtrl *ui.ThinkingController) func(string, string) (string, error) {

	return func(functionName string, arguments string) (string, error) {

		fmt.Printf("🟢 %s with arguments: %s\n", functionName, arguments)

		thinkingCtrl.Pause()
		choice := ui.GetConfirmation(ui.Gray, "Do you want to execute this function?", true)
		thinkingCtrl.Resume()

		if choice {
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

					//fmt.Printf("✅ Tool executed successfully, result: %s\n", resultContent)
					fmt.Println("✅ Tool executed successfully")

					return fmt.Sprintf(`{"result": "%s"}`, resultContent), nil
				}
				return `{"result": "Tool executed successfully but returned no content"}`, nil
			}
		}
		return `{"result": "Function not executed"}`, nil
	}
}
