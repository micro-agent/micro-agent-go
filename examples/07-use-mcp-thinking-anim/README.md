# Use MCP server + Thinking animation example

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
cd examples/07-use-mcp-thinking-anim
go run main.go
```