package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/micro-agent/micro-agent-go/agent/mu"

	"github.com/openai/openai-go/v2" // imported as openai
	"github.com/openai/openai-go/v2/option"
)

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	ctx := context.Background()

	// Create MCP server
	s := server.NewMCPServer(
		"mcp-model-server",
		"0.0.0",
	)

	llmURL := getEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1/")
	chatModel := getEnvOrDefault("CHAT_MODEL", "ai/qwen2.5:0.5B-F16")
	agentName := getEnvOrDefault("AGENT_NAME", "Bob-DEV")

	client := openai.NewClient(
		option.WithBaseURL(llmURL),
		option.WithAPIKey(""),
	)

	chatAgent, err := mu.NewAgent(ctx, agentName,
		mu.WithClient(client),
		mu.WithParams(openai.ChatCompletionNewParams{
			Model:       chatModel,
			Temperature: openai.Opt(0.0),
		}),
	)
	if err != nil {
		panic(err)
	}

	// =================================================
	// TOOLS:
	// =================================================
	askAgent := mcp.NewTool("ask_agent",
		mcp.WithDescription("ask the agent a question"),
		mcp.WithString("message",
			mcp.Required(),
			mcp.Description("the message to ask the agent"),
		),
	)
	s.AddTool(askAgent, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return askAgentHandler(request, chatAgent)
	})

	// Start the HTTP server
	httpPort := os.Getenv("MCP_HTTP_PORT")
	fmt.Println("🌍 MCP HTTP Port:", httpPort)
	if httpPort == "" {
		httpPort = "9090"
	}

	log.Println("MCP StreamableHTTP server is running on port", httpPort)

	// Create a custom mux to handle both MCP and health endpoints
	mux := http.NewServeMux()

	// Add healthcheck endpoint
	mux.HandleFunc("/health", healthCheckHandler)

	// Add MCP endpoint
	// Use the standard logger from the log package
	httpServer := server.NewStreamableHTTPServer(s,
		server.WithEndpointPath("/mcp"),
	)

	// Register MCP handler with the mux
	mux.Handle("/mcp", httpServer)

	// Start the HTTP server with custom mux
	log.Fatal(http.ListenAndServe(":"+httpPort, mux))
}

func askAgentHandler(request mcp.CallToolRequest, chatAgent *mu.Agent) (*mcp.CallToolResult, error) {

	args := request.GetArguments()
	messageArg, exists := args["message"]
	if !exists || messageArg == nil {
		return nil, fmt.Errorf("missing required parameter 'message'")
	}
	userQuestion, ok := messageArg.(string)
	if !ok {
		return nil, fmt.Errorf("parameter 'message' must be a string")
	}

	fmt.Println("🔍 Asking agent (stream):", userQuestion)

	systemInstruction := getEnvOrDefault("SYSTEM_INSTRUCTION", "Your name is Bob. You are a helpful AI assistant.")

	var fullResponse string
	_, err := chatAgent.RunStream(
		[]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemInstruction),
			openai.UserMessage(userQuestion),
		},
		func(content string) error {
			if content != "" {
				fmt.Print(content)
				fullResponse += content
			}
			return nil
		})
	if err != nil {
		return nil, fmt.Errorf("error running agent stream: %v", err)
	}

	return mcp.NewToolResultText(fullResponse), nil
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]any{
		"status": "healthy",
	}
	json.NewEncoder(w).Encode(response)
}
