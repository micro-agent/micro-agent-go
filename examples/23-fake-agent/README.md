# 23-fake-agent

This example demonstrates how to implement a fake/mock version of the `mu.Agent` interface. This is useful for testing, development, or creating alternative implementations that don't rely on external AI services.

## What this example shows

- **Interface Implementation**: How to create a custom agent that implements the `mu.Agent` interface
- **Mock Responses**: Simulating AI responses without making actual API calls
- **Streaming Simulation**: Mimicking real-time streaming behavior with delays
- **Tool Call Simulation**: Fake tool detection and execution
- **Embedding Generation**: Creating pseudo-embeddings based on content
- **Message Management**: Implementing GetMessages() and SetMessages()

## Key Features

The `FakeAgent` implements all methods required by the `mu.Agent` interface:

- `Run()` - Basic completion simulation
- `RunStream()` - Streaming completion with simulated delays
- `RunWithReasoning()` - Response with fake reasoning
- `RunStreamWithReasoning()` - Streaming with both content and reasoning
- `DetectToolCalls()` - Tool detection and execution simulation
- `DetectToolCallsStream()` - Streaming tool calls
- `GenerateEmbeddingVector()` - Pseudo-embedding generation
- `GetMessages()` / `SetMessages()` - Message management

## Benefits of this approach

1. **Testing**: Perfect for unit tests where you don't want real AI API calls
2. **Development**: Develop and test your application logic without API costs
3. **Offline Development**: Work without internet connection
4. **Predictable Responses**: Consistent behavior for testing scenarios
5. **Interface Demonstration**: Shows how the Agent interface enables pluggable implementations

## How to run

```bash
go run main.go
```

## Example Output

```
🤖 Fake Agent Example
=====================

1. Testing Run():
Response: Hello! I'm FakeBot, your fake AI assistant. How can I help you today?

2. Testing RunStream():
Streamed response: I'm a fake agent, so I can't check real weather, but let's pretend it's sunny and 72°F!
Full response: I'm a fake agent, so I can't check real weather, but let's pretend it's sunny and 72°F!

3. Testing RunWithReasoning():
Reasoning: 🤔 Reasoning: The user asked about 'What's 2+2?'. As a fake agent, I'm simulating the thought process...
Content: I calculated that 2+2 = 4 (even fake agents know basic math!)

4. Testing DetectToolCalls():
   🔧 Executing tool: search_tool with args: {"query": "example query"}
Finish reason: stop
Tool results: [Fake result from search_tool]
Assistant message: I found some results using the search tool: Fake result from search_tool

5. Testing GenerateEmbeddingVector():
Generated embedding vector with 1536 dimensions
First 5 values: -0.4370, -0.4280, -0.4310, -0.4360, -0.4470

6. Testing GetMessages() and SetMessages():
Current message count: 8
Updated message count: 9

✅ All fake agent methods tested successfully!
```
