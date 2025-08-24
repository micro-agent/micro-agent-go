package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/micro-agent/micro-agent-go/agent/helpers"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/shared"
	//"github.com/openai/openai-go/v2/shared/constant"
)

// MCPClient wraps an MCP client connection with available tools
type MCPClient struct {
	mcpclient   *client.Client
	ToolsResult *mcp.ListToolsResult
}

// NewStreamableHttpMCPClient creates and initializes a new MCP client over HTTP
func NewStreamableHttpMCPClient(ctx context.Context, mcpHostURL string) (*MCPClient, error) {
	mcpClient, err := client.NewStreamableHttpClient(
		mcpHostURL, // Use environment variable for MCP host
	)
	//defer mcpClient.Close()
	if err != nil {
		return nil, err
	}
	// Start the connection to the server
	err = mcpClient.Start(ctx)
	if err != nil {
		return nil, err
	}

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "micro agent",
		Version: "0.0.0",
	}
	_, err = mcpClient.Initialize(ctx, initRequest)
	if err != nil {
		return nil, err
	}
	//fmt.Println("Streamable HTTP client connected & initialized with server!", result)
	//ui.Println(ui.Yellow, "Streamable HTTP client connected & initialized with server!")

	toolsRequest := mcp.ListToolsRequest{}
	mcpTools, err := mcpClient.ListTools(ctx, toolsRequest)
	if err != nil {
		return nil, err
	}

	return &MCPClient{
		mcpclient:   mcpClient,
		ToolsResult: mcpTools,
	}, nil
}

// OpenAITools converts the MCP client's tools to OpenAI-compatible format
func (c *MCPClient) OpenAITools() []openai.ChatCompletionToolUnionParam {
	return ConvertMCPToolsToOpenAITools(c.ToolsResult)
}

// OpenAIToolsWithFilter converts only the filtered MCP tools to OpenAI-compatible format
func (c *MCPClient) OpenAIToolsWithFilter(toolsFilter []string) []openai.ChatCompletionToolUnionParam {
	return ConvertMCPToolsToOpenAIToolsWithFilter(c.ToolsResult, toolsFilter)
}

// Close safely closes the MCP client connection
func (c *MCPClient) Close() error {
	if c.mcpclient != nil {
		return c.mcpclient.Close()
	}
	return nil
}

// CallTool executes a tool call with the given function name and JSON arguments
func (c *MCPClient) CallTool(ctx context.Context, functionName string, arguments string) (*mcp.CallToolResult, error) {

	// Parse the tool arguments from JSON string
	var args map[string]any
	args, _ = helpers.JsonStringToMap(arguments)
	// TODO: check if this is useful for the request

	// NOTE: Call the MCP tool with the arguments
	request := mcp.CallToolRequest{}
	request.Params.Name = functionName
	request.Params.Arguments = args

	// NOTE: Call the tool using the MCP client
	toolResponse, err := c.mcpclient.CallTool(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error calling tool %s: %w", functionName, err)
	}
	if toolResponse == nil || len(toolResponse.Content) == 0 {
		return nil, fmt.Errorf("no content returned from tool %s", functionName)
	}

	return toolResponse, nil
}

// ConvertMCPToolsToOpenAITools transforms MCP tool definitions into OpenAI tool format
func ConvertMCPToolsToOpenAITools(tools *mcp.ListToolsResult) []openai.ChatCompletionToolUnionParam {
	openAITools := make([]openai.ChatCompletionToolUnionParam, len(tools.Tools))
	for i, tool := range tools.Tools {

		openAITools[i] = openai.ChatCompletionFunctionTool(shared.FunctionDefinitionParam{
			Name:        tool.Name,
			Description: openai.String(tool.Description),
			Parameters: shared.FunctionParameters{
				"type":       "object",
				"properties": tool.InputSchema.Properties,
				"required":   tool.InputSchema.Required,
			},
		},
		)
	}
	return openAITools
}

// ConvertMCPToolsToOpenAIToolsWithFilter transforms filtered MCP tool definitions into OpenAI tool format
func ConvertMCPToolsToOpenAIToolsWithFilter(tools *mcp.ListToolsResult, toolsFilter []string) []openai.ChatCompletionToolUnionParam {
	// Create a set for quick lookup of allowed tool names
	allowedTools := make(map[string]bool)
	for _, name := range toolsFilter {
		allowedTools[name] = true
	}

	// Filter tools and convert to OpenAI format
	var openAITools []openai.ChatCompletionToolUnionParam
	for _, tool := range tools.Tools {
		if allowedTools[tool.Name] {
			openAITools = append(openAITools, openai.ChatCompletionFunctionTool(shared.FunctionDefinitionParam{
				Name:        tool.Name,
				Description: openai.String(tool.Description),
				Parameters: shared.FunctionParameters{
					"type":       "object",
					"properties": tool.InputSchema.Properties,
					"required":   tool.InputSchema.Required,
				},
			},
			))
		}
	}
	return openAITools
}
