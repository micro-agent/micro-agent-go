package mu

import (
	"errors"
	"fmt"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/shared/constant"
)

// DetectToolCalls processes a conversation with tool calls support.
// It handles the complete tool calling workflow: detecting tool calls, executing them via callback,
// and managing the conversation history until completion.
//
// Parameters:
//   - messages: Initial conversation messages to start with
//   - callBack: Function to execute when tools are called. Takes functionName and arguments (JSON string),
//     returns the result as a JSON string
//
// Returns:
//   - finishReason: The reason the conversation ended ("stop" for normal completion, other values for errors)
//   - results: Slice of all tool execution results (JSON strings)
//   - lastAssistantMessage: The final message from the assistant when conversation ends normally
//   - error: Any error that occurred during processing
func (agent *BasicAgent) DetectToolCalls(messages []openai.ChatCompletionMessageParamUnion, toolCallBack func(functionName string, arguments string) (string, error)) (string, []string, string, error) {

	stopped := false
	results := []string{}
	lastAssistantMessage := ""
	finishReason := ""

	for !stopped {
		// TOOL: Make a function call request
		//fmt.Println("⏳ Making function call request...")

		agent.Params.Messages = messages

		completion, err := agent.Client.Chat.Completions.New(agent.ctx, agent.Params)
		if err != nil {
			return "", results, "", err
			//return nil, errors.New("error making function call request [completion]")
		}

		finishReason = completion.Choices[0].FinishReason

		// Extract reasoning_content from RawJSON
		// completion.Choices[0].Message.RawJSON()

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

				// Create assistant message with tool calls using proper union type
				assistantMessage := openai.ChatCompletionMessageParamUnion{
					OfAssistant: &openai.ChatCompletionAssistantMessageParam{
						ToolCalls: toolCallParams,
					},
				}

				// Add the assistant message with tool calls to the conversation history
				messages = append(messages, assistantMessage)

				// TOOL: Process each detected tool call
				//fmt.Println("🚀 Processing tool calls...")

				for _, toolCall := range detectedToolCalls {
					functionName := toolCall.Function.Name
					functionArgs := toolCall.Function.Arguments
					//callID := toolCall.ID

					// TOOL: Execute the function with the provided arguments
					//fmt.Printf("▶️ Executing function: %s with args: %s\n", functionName, functionArgs)

					resultContent, errExec := toolCallBack(functionName, functionArgs)

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

					//fmt.Printf("Function result: %s with CallID: %s\n\n", resultContent, callID)

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
				// TODO: Handle case where no tool calls were detected
				fmt.Println("😢 No tool calls found in response")
			}

		case "stop":
			//fmt.Println("🟥 Stopping due to 'stop' finish reason.")
			stopped = true
			lastAssistantMessage = completion.Choices[0].Message.Content
			//fmt.Printf("🤖 %s\n", lastAssistantMessage)

			// Add final assistant message to conversation history
			messages = append(messages, openai.AssistantMessage(lastAssistantMessage))

		default:
			//fmt.Printf("🔴 Unexpected response: %s\n", finishReason)
			stopped = true
		}

	}
	return finishReason, results, lastAssistantMessage, nil
}
