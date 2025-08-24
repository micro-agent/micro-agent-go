# Tool streaming completion example

## Pre-requisites

- Install Docker Model Runner
- Pull the model image:
  ```bash
  docker model pull hf.co/menlo/jan-nano-gguf:q4_k_m
  ```

## Running the Example

```bash
cd examples/04-tool-stream-completion
go run main.go
```