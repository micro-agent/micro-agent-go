package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/micro-agent/micro-agent-go/agent/mu"
	"github.com/openai/openai-go/v2"
)

// FakeAgent is a mock implementation of the mu.Agent interface
// It simulates AI responses without making actual API calls
type FakeAgent struct {
	name     string
	messages []openai.ChatCompletionMessageParamUnion
}

// NewFakeAgent creates a new fake agent instance
func NewFakeAgent(name string) mu.Agent {
	return &FakeAgent{
		name:     name,
		messages: []openai.ChatCompletionMessageParamUnion{},
	}
}

// Run simulates a simple completion
func (f *FakeAgent) Run(Messages []openai.ChatCompletionMessageParamUnion) (string, error) {
	f.messages = append(f.messages, Messages...)
	
	// Extract user message content for simulation
	var userMessage string
	for _, msg := range Messages {
		if msg.OfUser != nil {
			if msg.OfUser.Content.OfString.Value != "" {
				userMessage = msg.OfUser.Content.OfString.Value
				break
			}
		}
	}
	
	return f.simulateResponse(userMessage), nil
}

// RunStream simulates streaming completion
func (f *FakeAgent) RunStream(Messages []openai.ChatCompletionMessageParamUnion, callBack func(content string) error) (string, error) {
	f.messages = append(f.messages, Messages...)
	
	// Extract user message content for simulation
	var userMessage string
	for _, msg := range Messages {
		if msg.OfUser != nil {
			if msg.OfUser.Content.OfString.Value != "" {
				userMessage = msg.OfUser.Content.OfString.Value
				break
			}
		}
	}
	
	response := f.simulateResponse(userMessage)
	
	// Simulate streaming by sending chunks
	words := strings.Fields(response)
	fullResponse := ""
	
	for _, word := range words {
		chunk := word + " "
		fullResponse += chunk
		
		// Simulate streaming delay
		time.Sleep(50 * time.Millisecond)
		
		if err := callBack(chunk); err != nil {
			return fullResponse, err
		}
	}
	
	return fullResponse, nil
}

// RunWithReasoning simulates completion with reasoning
func (f *FakeAgent) RunWithReasoning(Messages []openai.ChatCompletionMessageParamUnion) (string, string, error) {
	f.messages = append(f.messages, Messages...)
	
	// Extract user message content for simulation
	var userMessage string
	for _, msg := range Messages {
		if msg.OfUser != nil {
			if msg.OfUser.Content.OfString.Value != "" {
				userMessage = msg.OfUser.Content.OfString.Value
				break
			}
		}
	}
	
	response := f.simulateResponse(userMessage)
	reasoning := f.simulateReasoning(userMessage)
	
	return response, reasoning, nil
}

// RunStreamWithReasoning simulates streaming completion with reasoning
func (f *FakeAgent) RunStreamWithReasoning(Messages []openai.ChatCompletionMessageParamUnion, contentCallback func(content string) error, reasoningCallback func(reasoning string) error) (string, string, error) {
	f.messages = append(f.messages, Messages...)
	
	// Extract user message content for simulation
	var userMessage string
	for _, msg := range Messages {
		if msg.OfUser != nil {
			if msg.OfUser.Content.OfString.Value != "" {
				userMessage = msg.OfUser.Content.OfString.Value
				break
			}
		}
	}
	
	response := f.simulateResponse(userMessage)
	reasoning := f.simulateReasoning(userMessage)
	
	// Stream reasoning first
	reasoningWords := strings.Fields(reasoning)
	fullReasoning := ""
	for _, word := range reasoningWords {
		chunk := word + " "
		fullReasoning += chunk
		time.Sleep(30 * time.Millisecond)
		
		if err := reasoningCallback(chunk); err != nil {
			return "", fullReasoning, err
		}
	}
	
	// Then stream content
	responseWords := strings.Fields(response)
	fullResponse := ""
	for _, word := range responseWords {
		chunk := word + " "
		fullResponse += chunk
		time.Sleep(50 * time.Millisecond)
		
		if err := contentCallback(chunk); err != nil {
			return fullResponse, fullReasoning, err
		}
	}
	
	return fullResponse, fullReasoning, nil
}

// DetectToolCalls simulates tool detection and execution
func (f *FakeAgent) DetectToolCalls(messages []openai.ChatCompletionMessageParamUnion, toolCallBack func(functionName string, arguments string) (string, error)) (string, []string, string, error) {
	f.messages = append(f.messages, messages...)
	
	// Extract user message content
	var userMessage string
	for _, msg := range messages {
		if msg.OfUser != nil {
			if msg.OfUser.Content.OfString.Value != "" {
				userMessage = msg.OfUser.Content.OfString.Value
				break
			}
		}
	}
	
	// Simulate tool detection
	results := []string{}
	assistantMessage := ""
	
	if strings.Contains(strings.ToLower(userMessage), "search") {
		// Simulate a search tool call
		result, err := toolCallBack("search_tool", `{"query": "example query"}`)
		if err != nil {
			return "error", results, "", err
		}
		results = append(results, result)
		assistantMessage = "I found some results using the search tool: " + result
	} else if strings.Contains(strings.ToLower(userMessage), "calculate") {
		// Simulate a calculator tool call
		result, err := toolCallBack("calculator", `{"expression": "2+2"}`)
		if err != nil {
			return "error", results, "", err
		}
		results = append(results, result)
		assistantMessage = "I calculated the result: " + result
	} else {
		// No tools needed
		assistantMessage = f.simulateResponse(userMessage)
	}
	
	return "stop", results, assistantMessage, nil
}

// DetectToolCallsStream simulates streaming tool detection and execution
func (f *FakeAgent) DetectToolCallsStream(messages []openai.ChatCompletionMessageParamUnion, toolCallback func(functionName string, arguments string) (string, error), streamCallback func(content string) error) (string, []string, string, error) {
	f.messages = append(f.messages, messages...)
	
	// Extract user message content
	var userMessage string
	for _, msg := range messages {
		if msg.OfUser != nil {
			if msg.OfUser.Content.OfString.Value != "" {
				userMessage = msg.OfUser.Content.OfString.Value
				break
			}
		}
	}
	
	results := []string{}
	assistantMessage := ""
	
	if strings.Contains(strings.ToLower(userMessage), "search") {
		// Stream thinking process
		thinking := "Let me search for that information..."
		for _, char := range thinking {
			time.Sleep(20 * time.Millisecond)
			if err := streamCallback(string(char)); err != nil {
				return "error", results, "", err
			}
		}
		
		result, err := toolCallback("search_tool", `{"query": "example query"}`)
		if err != nil {
			return "error", results, "", err
		}
		results = append(results, result)
		assistantMessage = "I found some results using the search tool: " + result
		
		// Stream final result
		if err := streamCallback("\n" + assistantMessage); err != nil {
			return "error", results, assistantMessage, err
		}
	} else {
		// No tools needed, just stream response
		response := f.simulateResponse(userMessage)
		assistantMessage = response
		
		words := strings.Fields(response)
		for _, word := range words {
			chunk := word + " "
			time.Sleep(50 * time.Millisecond)
			if err := streamCallback(chunk); err != nil {
				return "error", results, assistantMessage, err
			}
		}
	}
	
	return "stop", results, assistantMessage, nil
}

// GenerateEmbeddingVector simulates embedding generation
func (f *FakeAgent) GenerateEmbeddingVector(content string) ([]float64, error) {
	// Generate a fake embedding vector based on content length and characters
	vector := make([]float64, 1536) // Standard OpenAI embedding size
	
	for i := range vector {
		// Create pseudo-random values based on content and position
		seed := float64(len(content) + int(content[i%len(content)]) + i)
		vector[i] = (seed * 0.001) - 0.5 // Values between -0.5 and 0.5
	}
	
	return vector, nil
}

// GetMessages returns the current messages
func (f *FakeAgent) GetMessages() []openai.ChatCompletionMessageParamUnion {
	return f.messages
}

// SetMessages sets the messages
func (f *FakeAgent) SetMessages(messages []openai.ChatCompletionMessageParamUnion) {
	f.messages = messages
}

// simulateResponse generates a fake AI response based on the input
func (f *FakeAgent) simulateResponse(userMessage string) string {
	responses := map[string]string{
		"hello":     fmt.Sprintf("Hello! I'm %s, your fake AI assistant. How can I help you today?", f.name),
		"weather":   "I'm a fake agent, so I can't check real weather, but let's pretend it's sunny and 72°F!",
		"code":      "Here's some fake code: `func main() { fmt.Println(\"Hello from fake agent!\") }`",
		"time":      "The current time is... well, I'm fake, so let's say it's always coffee time! ☕",
		"calculate": "I calculated that 2+2 = 4 (even fake agents know basic math!)",
		"search":    "I found exactly what you were looking for! (Just kidding, I'm a fake agent)",
	}
	
	userLower := strings.ToLower(userMessage)
	for keyword, response := range responses {
		if strings.Contains(userLower, keyword) {
			return response
		}
	}
	
	return fmt.Sprintf("I'm %s, a fake AI agent. You said: \"%s\". I don't have real AI capabilities, but I'm pretending to understand and respond!", f.name, userMessage)
}

// simulateReasoning generates fake reasoning content
func (f *FakeAgent) simulateReasoning(userMessage string) string {
	return fmt.Sprintf("🤔 Reasoning: The user asked about '%s'. As a fake agent, I'm simulating the thought process of analyzing this request and determining the best response approach.", userMessage)
}

func main() {
	fmt.Println("🤖 Fake Agent Example")
	fmt.Println("=====================")
	
	// Create a fake agent
	fakeAgent := NewFakeAgent("FakeBot")
	
	// Test basic run
	fmt.Println("\n1. Testing Run():")
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage("Hello, how are you?"),
	}
	
	response, err := fakeAgent.Run(messages)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response: %s\n", response)
	
	// Test streaming
	fmt.Println("\n2. Testing RunStream():")
	streamMessages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage("Tell me about the weather"),
	}
	
	fmt.Print("Streamed response: ")
	streamResponse, err := fakeAgent.RunStream(streamMessages, func(content string) error {
		fmt.Print(content)
		return nil
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nFull response: %s\n", streamResponse)
	
	// Test reasoning
	fmt.Println("\n3. Testing RunWithReasoning():")
	reasoningMessages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage("What's 2+2?"),
	}
	
	content, reasoning, err := fakeAgent.RunWithReasoning(reasoningMessages)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Reasoning: %s\n", reasoning)
	fmt.Printf("Content: %s\n", content)
	
	// Test tool calls
	fmt.Println("\n4. Testing DetectToolCalls():")
	toolMessages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage("Please search for information about Go programming"),
	}
	
	toolExecutor := func(functionName string, arguments string) (string, error) {
		fmt.Printf("   🔧 Executing tool: %s with args: %s\n", functionName, arguments)
		return fmt.Sprintf("Fake result from %s", functionName), nil
	}
	
	finishReason, results, assistantMsg, err := fakeAgent.DetectToolCalls(toolMessages, toolExecutor)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Finish reason: %s\n", finishReason)
	fmt.Printf("Tool results: %v\n", results)
	fmt.Printf("Assistant message: %s\n", assistantMsg)
	
	// Test embeddings
	fmt.Println("\n5. Testing GenerateEmbeddingVector():")
	vector, err := fakeAgent.GenerateEmbeddingVector("Hello world")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Generated embedding vector with %d dimensions\n", len(vector))
	fmt.Printf("First 5 values: %.4f, %.4f, %.4f, %.4f, %.4f\n", 
		vector[0], vector[1], vector[2], vector[3], vector[4])
	
	// Test message management
	fmt.Println("\n6. Testing GetMessages() and SetMessages():")
	currentMessages := fakeAgent.GetMessages()
	fmt.Printf("Current message count: %d\n", len(currentMessages))
	
	// Add a new message
	newMessages := append(currentMessages, openai.SystemMessage("You are a helpful assistant"))
	fakeAgent.SetMessages(newMessages)
	
	updatedMessages := fakeAgent.GetMessages()
	fmt.Printf("Updated message count: %d\n", len(updatedMessages))
	
	fmt.Println("\n✅ All fake agent methods tested successfully!")
	fmt.Println("This demonstrates how the Agent interface can be implemented")
	fmt.Println("with different backends - real AI services or fake/mock implementations.")
}