package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/micro-agent/micro-agent-go/agent/mu"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"github.com/openai/openai-go/v2/shared"
)

func main() {

	ctx := context.Background()

	client := openai.NewClient(
		option.WithBaseURL("http://localhost:12434/engines/llama.cpp/v1"),
		option.WithAPIKey(""),
	)

	toolAgent, err := mu.NewAgent(ctx, "Bob",
		mu.WithClient(client),
		mu.WithParams(openai.ChatCompletionNewParams{
			Model:       "hf.co/menlo/jan-nano-gguf:q4_k_m",
			Temperature: openai.Opt(0.0),
			ToolChoice: openai.ChatCompletionToolChoiceOptionUnionParam{
				OfAuto: openai.String("auto"),
			},
			Tools:             GetToolsIndex(),
			ParallelToolCalls: openai.Opt(false),
		}),
	)
	if err != nil {
		panic(err)
	}

	// Say "Exit" to stop the process
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage(`
			Make the sum of 40 and 2,
			then say hello to Bob and to Sam,
			make the sum of 5 and 37
			Say hello to Alice			
		`),
	}

	finishReason, results, assistantMessage, err := toolAgent.DetectToolCalls(messages, executeFunction)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Finish Reason: %s\n", finishReason)
	fmt.Printf("Results: %v\n", results)
	fmt.Printf("Assistant Message: %s\n", assistantMessage)

}

func GetToolsIndex() []openai.ChatCompletionToolUnionParam {
	calculateSumTool := openai.ChatCompletionFunctionTool(shared.FunctionDefinitionParam{
		Name:        "calculate_sum",
		Description: openai.String("Calculate the sum of two numbers"),
		Parameters: shared.FunctionParameters{
			"type": "object",
			"properties": map[string]interface{}{
				"a": map[string]string{
					"type":        "number",
					"description": "The first number",
				},
				"b": map[string]string{
					"type":        "number",
					"description": "The second number",
				},
			},
			"required": []string{"a", "b"},
		},
	})

	sayHelloTool := openai.ChatCompletionFunctionTool(shared.FunctionDefinitionParam{
		Name:        "say_hello",
		Description: openai.String("Say hello to the given name"),
		Parameters: shared.FunctionParameters{
			"type": "object",
			"properties": map[string]interface{}{
				"name": map[string]string{
					"type":        "string",
					"description": "The name to greet",
				},
			},
			"required": []string{"name"},
		},
	})

	sayExit := openai.ChatCompletionFunctionTool(shared.FunctionDefinitionParam{
		Name:        "say_exit",
		Description: openai.String("Say exit"),
		Parameters: shared.FunctionParameters{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
	})

	return []openai.ChatCompletionToolUnionParam{
		calculateSumTool,
		sayHelloTool,
		sayExit,
	}
}

func executeFunction(functionName string, arguments string) (string, error) {
	fmt.Printf("🟢 Executing function: %s with arguments: %s\n", functionName, arguments)
	switch functionName {
	case "say_hello":
		var args struct {
			Name string `json:"name"`
		}
		if err := json.Unmarshal([]byte(arguments), &args); err != nil {
			return `{"error": "Invalid arguments for say_hello"}`, nil
		}
		hello := fmt.Sprintf("👋 Hello, %s!🙂", args.Name)
		return fmt.Sprintf(`{"message": "%s"}`, hello), nil

	case "calculate_sum":
		var args struct {
			A float64 `json:"a"`
			B float64 `json:"b"`
		}
		if err := json.Unmarshal([]byte(arguments), &args); err != nil {
			return `{"error": "Invalid arguments for calculate_sum"}`, nil
		}
		sum := args.A + args.B
		return fmt.Sprintf(`{"result": %g}`, sum), nil

	case "say_exit":

		// NOTE: Returning a message and an ExitToolCallsLoopError to stop further processing
		return fmt.Sprintf(`{"message": "%s"}`, "❌ EXIT"), &mu.ExitToolCallsLoopError{Message: "❌ EXIT"}

	default:
		return `{"error": "Unknown function"}`, &mu.ExitToolCallsLoopError{Message: fmt.Sprintf("Unknown function: %s", functionName)}
	}
}
