# Streaming completion example

## Pre-requisites

- Install Docker Model Runner
- Pull the model image:
  ```bash
  docker model pull ai/qwen2.5:1.5B-F16
  ```

## Running the Example

```bash
cd examples/02-stream-completion
go run main.go
```