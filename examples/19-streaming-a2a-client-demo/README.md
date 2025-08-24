# Streaming A2A Client Demo

This example demonstrates how to use the A2A (Agent-to-Agent) client to communicate with a µ-agent server using streaming responses.

## Overview

The demo showcases:
- Connecting to a µ-agent server
- Pinging the agent to verify connectivity and retrieve agent metadata
- Sending a streaming task request
- Processing streamed response chunks in real-time
- Handling the final response

## Usage

1. Start a µ-agent server on `http://localhost:7777`
2. Run the demo:
   ```bash
   go run main.go
   ```

## Features

- **Connection Verification**: Pings the agent and displays its metadata
- **Streaming Communication**: Sends requests and receives responses in real-time chunks
- **Real-time Display**: Shows content as it arrives from the agent
- **Response Summary**: Displays final task status and response statistics

## Example Output

```
🔍 Pinging agent...
✅ Connected to agent: My Agent
📝 Description: A helpful AI agent
🔧 Available skills: 5

🚀 Sending streaming task request: task-1692123456
🌊 Streaming response:
================================================================================
Once upon a time, there was a curious robot named Chef-Bot who dreamed of...
[content streams here in real-time]
================================================================================
✅ Task completed: task-1692123456
🎯 Status: completed
📝 Total response length: 1247 characters

🏁 Streaming demo completed!
```

## Notes

- This uses the experimental A2A package which is subject to change
- The agent must support streaming responses via the `/stream` endpoint
- The example uses Server-Sent Events (SSE) format for streaming