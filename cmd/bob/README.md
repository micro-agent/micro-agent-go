# Bob, µ-agent CLI

This directory contains a complete command-line application that demonstrates some capabilities of the micro-agent library.

## Interactive AI Assistant

The `cmd/bob/main.go` program demonstrates a more advanced implementation of µAgent, featuring an interactive AI assistant with MCP (Model Context Protocol) tool integration.

### Architecture Overview

```mermaid
graph TD
    A[User Input] --> B[Bob CLI Application]
    B --> C[Environment Configuration]
    C --> D[OpenAI Client Setup]
    C --> E[MCP Client Setup]
    D --> F[µAgent Creation]
    E --> F
    F --> G[Interactive Loop]
    G --> H[User Prompt]
    H --> I{Input = /bye?}
    I -->|Yes| J[Exit Application]
    I -->|No| K[Process Message]
    K --> L[Tool Detection & Execution]
    L --> M[Stream Response]
    M --> N[Display Result]
    N --> G
```

### Component Architecture

```mermaid
graph LR
    subgraph "Bob Application"
        A[Main Function] --> B[Environment Config]
        A --> C[OpenAI Client]
        A --> D[MCP Client]
        A --> E[µAgent]
        A --> F[Interactive Loop]
    end
    
    subgraph "External Services"
        G[LLM Provider]
        H[MCP Tools Server]
    end
    
    subgraph "User Interface"
        I[CLI Prompt]
        J[Markdown Display]
        K[Thinking Controller]
    end
    
    C --> G
    D --> H
    F --> I
    F --> J
    F --> K
```

### Program Flow

```mermaid
sequenceDiagram
    participant User
    participant Bob as Bob Application
    participant LLM as LLM Provider
    participant MCP as MCP Tools Server
    participant UI as User Interface

    User->>Bob: Start Application
    Bob->>Bob: Load Environment Variables
    Bob->>LLM: Initialize OpenAI Client
    Bob->>MCP: Initialize MCP Client
    Bob->>MCP: Fetch Available Tools
    Bob->>Bob: Create µAgent with Tools
    
    loop Interactive Session
        Bob->>UI: Display Prompt
        User->>Bob: Enter Command
        
        alt Command is "/bye"
            Bob->>User: Display "Goodbye!"
            Bob->>Bob: Exit Loop
        else Process Command
            Bob->>UI: Start Thinking Animation
            Bob->>LLM: Send Message with Tools
            LLM->>Bob: Return Tool Calls (if any)
            
            opt Tool Execution Required
                Bob->>UI: Pause Animation & Ask Confirmation
                User->>Bob: Approve/Reject/Abort
                
                alt Approved
                    Bob->>MCP: Execute Tool
                    MCP->>Bob: Return Tool Result
                    Bob->>LLM: Send Tool Result
                else Rejected/Aborted
                    Bob->>Bob: Skip Tool Execution
                end
            end
            
            Bob->>UI: Start Streaming Animation
            LLM->>Bob: Stream Response Content
            Bob->>UI: Display Markdown Response
        end
    end
```

### Configuration

The application supports configuration through environment variables:

| Variable | Default Value | Description |
|----------|---------------|-------------|
| `PROVIDER_BASE_URL` | `http://localhost:12434/engines/llama.cpp/v1` | LLM provider base URL |
| `PROVIDER_API_KEY` | `""` | API key for the LLM provider |
| `MCP_HOST_URL` | `http://localhost:9011` | MCP tools server URL |
| `MODEL_ID` | `hf.co/menlo/jan-nano-gguf:q4_k_m` | Model identifier |
| `SYSTEM_MESSAGE` | Bob the Bot default message | System prompt for the AI assistant |

### Key Features

1. **Environment-based Configuration**: All connection parameters are configurable via environment variables with sensible defaults.

2. **MCP Tool Integration**: Connects to an MCP (Model Context Protocol) server to provide the AI with external tools and capabilities.

3. **Interactive User Confirmation**: Before executing any tool, the application asks for user confirmation with options to approve, reject, or abort.

4. **Real-time Streaming**: Responses are streamed in real-time with visual feedback through thinking and streaming controllers.

5. **Graceful Exit**: Users can type `/bye` to exit the application cleanly.

6. **Rich UI Feedback**: Uses colored output and animations to provide clear visual feedback during different stages of processing.

### Tool Execution Flow

```mermaid
stateDiagram-v2
    [*] --> DetectTools: Message with Tools
    DetectTools --> ToolFound: Tools Detected
    DetectTools --> NoTools: No Tools
    
    ToolFound --> PromptUser: Display Tool & Arguments
    PromptUser --> UserApproves: User selects "y"
    PromptUser --> UserRejects: User selects "n" 
    PromptUser --> UserAborts: User selects "a"
    
    UserApproves --> ExecuteTool: Call MCP Server
    ExecuteTool --> ToolResult: Tool Executed
    ToolResult --> StreamResponse: Send Result to LLM
    
    UserRejects --> SkipTool: Return "Not executed"
    SkipTool --> StreamResponse
    
    UserAborts --> ExitLoop: Abort Tool Loop
    
    NoTools --> StreamResponse: Direct LLM Response
    StreamResponse --> DisplayResult: Show Markdown
    DisplayResult --> [*]
    
    ExitLoop --> [*]
```

### Prerequisites

**Start MCP Server**: Before running the Bob application, you must start an MCP (Model Context Protocol) server. The application expects an MCP server to be running at `http://localhost:9011` by default. This server provides external tools and capabilities that extend the AI's functionality beyond basic conversation.

Without a running MCP server, the application will fail to initialize properly since MCP tool integration is a core feature.

### Running the Application

```bash
# With default configuration
go run cmd/bob/main.go

# With custom configuration
PROVIDER_BASE_URL="https://api.openai.com/v1" \
PROVIDER_API_KEY="your-api-key" \
MODEL_ID="gpt-4" \
go run cmd/bob/main.go
```