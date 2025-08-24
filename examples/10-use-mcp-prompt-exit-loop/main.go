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
	toolsIndex := mcpClient.OpenAITools()
	for _, tool := range toolsIndex {
		ui.Printf(ui.Magenta, "Tool: %s - %s\n", tool.GetFunction().Name, tool.GetFunction().Description)
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

	content, _ := ui.SimplePrompt("🤖 (/bye to exit)>", "Type your command here...")

	// Say "Exit" to stop the process
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(`
			Your name is "Bob the Bot".
			You are a helpful assistant that can search for code snippets and generate markdown documentation.
			You are an expert in Python programming and can explain code snippets in detail.
		`),
		openai.UserMessage(content.Input),
	}
	/*
		search a snippet for doing k33g program in Python,
		then generate a markdow document from the found snippet
		and explain it	and sign with your name
	*/

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

// executeFunction creates a function that executes MCP tools with user confirmation.
// It returns a closure that prompts the user before executing any function call,
// allowing them to approve (y), reject (n), or abort the entire execution loop (a).
//
// Parameters:
//   - mcpClient: the MCP client for executing tools, can be nil
//   - thinkingCtrl: controller for pausing/resuming thinking animations during prompts
//
// Returns:
//   - func(string, string) (string, error): a function that takes functionName and arguments,
//     prompts for confirmation, and executes the tool if approved
func executeFunction(mcpClient *tools.MCPClient, thinkingCtrl *ui.ThinkingController) func(string, string) (string, error) {

	return func(functionName string, arguments string) (string, error) {

		fmt.Printf("🟢 %s with arguments: %s\n", functionName, arguments)

		thinkingCtrl.Pause()
		//choice := ui.GetConfirmation(ui.Gray, "Do you want to execute this function?", true)
		choice := ui.GetChoice(ui.Gray, "Do you want to execute this function? (y)es (n)o (a)bort", []string{"y",
			"n", "a"}, "y")
		thinkingCtrl.Resume()

		switch choice {
		case "n":
			return `{"result": "Function not executed"}`, nil
		case "a": // abort
			return `{"result": "Function not executed"}`,
				&mu.ExitToolCallsLoopError{Message: "Tool execution aborted by user"}

		default:

			// If MCP client is available, use it to execute the tool
			if mcpClient != nil {
				ctx := context.Background()
				result, err := mcpClient.CallTool(ctx, functionName, arguments)
				if err != nil {
					return "", fmt.Errorf("MCP tool execution failed: %v", err)
				}
				// resultContent = toolResponse.Content[0].(mcp.TextContent).Text
				// Convert MCP result to JSON string
				if len(result.Content) > 0 {
					// Take the first content item and return its text
					resultContent := result.Content[0].(mcp.TextContent).Text
					//fmt.Printf("✅ Tool executed successfully, result: %s\n", resultContent)
					fmt.Println("✅ Tool executed successfully")
					return fmt.Sprintf(`{"result": "%s"}`, resultContent), nil
				}
				return `{"result": "Tool executed successfully but returned no content"}`, nil
			}
			return `{"result": "Function not executed"}`, nil
		}
	}
}
