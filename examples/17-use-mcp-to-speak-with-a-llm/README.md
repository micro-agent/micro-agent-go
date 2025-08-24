# Use MCP server to speak with a LLM example

## Pre-requisites

- Install Docker Model Runner
- Pull the model image:
  ```bash
  docker model pull hf.co/menlo/jan-nano-gguf:q4_k_m
  ```
- Start the MCP servers:
  ```bash
  docker compose -f ../mcp-servers/compose.yml up 
  ```

## Running the Example

```bash
cd examples/17-use-mcp-to-speak-with-a-llm
go run main.go
```