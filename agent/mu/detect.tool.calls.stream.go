package mu

import (
	"errors"
	"fmt"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/shared/constant"
)

// DetectToolCallsStream processes a conversation with tool calls support using streaming.
// It handles the complete tool calling workflow with real-time streaming of assistant responses,
// detecting tool calls, executing them via callback, and managing the conversation history until completion.
//
// Parameters:
//   - messages: Initial conversation messages to start with
//   - streamCallback: Function called for each streaming chunk (content string) -> error
//   - toolCallback: Function to execute when tools are called. Takes functionName and arguments (JSON string),
//     returns the result as a JSON string
//
// Returns:
//   - finishReason: The reason the conversation ended ("stop" for normal completion, other values for errors)
//   - results: Slice of all tool execution results (JSON strings)
//   - lastAssistantMessage: The final message from the assistant when conversation ends normally
//   - error: Any error that occurred during processing
func (agent *Agent) DetectToolCallsStream(messages []openai.ChatCompletionMessageParamUnion, toolCallback func(functionName string, arguments string) (string, error), streamCallback func(content string) error) (string, []string, string, error) {
	stopped := false
	results := []string{}
	lastAssistantMessage := ""
	finishReason := ""

	for !stopped {
		agent.Params.Messages = messages

		stream := agent.Client.Chat.Completions.NewStreaming(agent.ctx, agent.Params)
		var response string
		var cbkRes error

		for stream.Next() {
			chunk := stream.Current()
			// Stream each chunk as it arrives
			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
				cbkRes = streamCallback(chunk.Choices[0].Delta.Content)
				response += chunk.Choices[0].Delta.Content
			}

			// if cbkRes != nil {
			// 	break
			// }

			if cbkRes != nil {
				var exitErr *ExitStreamCompletionError
				if errors.As(cbkRes, &exitErr) {
					break
				}
			}

		}

		if cbkRes != nil {
			return "", results, "", cbkRes
		}
		if err := stream.Err(); err != nil {
			return "", results, "", err
		}
		if err := stream.Close(); err != nil {
			return "", results, "", err
		}

		// Make a non-streaming call to get tool calls (streaming doesn't provide tool calls properly)
		completion, err := agent.Client.Chat.Completions.New(agent.ctx, agent.Params)
		if err != nil {
			return "", results, "", err
		}

		finishReason = completion.Choices[0].FinishReason

		switch finishReason {
		case "tool_calls":
			detectedToolCalls := completion.Choices[0].Message.ToolCalls

			if len(detectedToolCalls) > 0 {
				toolCallParams := make([]openai.ChatCompletionMessageToolCallUnionParam, len(detectedToolCalls))
				for i, toolCall := range detectedToolCalls {
					toolCallParams[i] = openai.ChatCompletionMessageToolCallUnionParam{
						OfFunction: &openai.ChatCompletionMessageFunctionToolCallParam{
							ID:   toolCall.ID,
							Type: constant.Function("function"),
							Function: openai.ChatCompletionMessageFunctionToolCallFunctionParam{
								Name:      toolCall.Function.Name,
								Arguments: toolCall.Function.Arguments,
							},
						},
					}
				}

				// Create assistant message with tool calls
				assistantMessage := openai.ChatCompletionMessageParamUnion{
					OfAssistant: &openai.ChatCompletionAssistantMessageParam{
						ToolCalls: toolCallParams,
					},
				}

				messages = append(messages, assistantMessage)

				// Execute each tool call
				for _, toolCall := range detectedToolCalls {
					functionName := toolCall.Function.Name
					functionArgs := toolCall.Function.Arguments

					resultContent, errExec := toolCallback(functionName, functionArgs)

					if errExec != nil {
						//fmt.Printf("🔴 Error executing function %s: %s\n", functionName, errExec.Error())
						var exitErr *ExitToolCallsLoopError
						if errors.As(errExec, &exitErr) {
							// If the error is an ExitLoopError, we stop processing further tool calls
							stopped = true
							finishReason = "exit_loop"
						} else {
							resultContent = fmt.Sprintf(`{"error": "Function execution failed: %s"}`, errExec.Error())
						}
					}

					if resultContent == "" {
						resultContent = `{"error": "Function execution returned empty result"}`
					}
					results = append(results, resultContent)

					// Add the tool call result to the conversation history
					messages = append(
						messages,
						openai.ToolMessage(
							resultContent,
							toolCall.ID,
						),
					)
				}

			} else {
				fmt.Println("😢 No tool calls found in response")
			}

		case "stop":
			stopped = true
			lastAssistantMessage = response

			// Add final assistant message to conversation history
			messages = append(messages, openai.AssistantMessage(lastAssistantMessage))

		default:
			stopped = true
		}
	}
	return finishReason, results, lastAssistantMessage, nil
}
