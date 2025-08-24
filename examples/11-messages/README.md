# Use MCP server + read the list of messages example
> 🚧 WIP
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
cd examples/11-messages
go run main.go
```