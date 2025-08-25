# Structured Output example

## Pre-requisites

- Install Docker Model Runner
- Pull the model image:
  ```bash
  docker model pull ai/qwen2.5:1.5B-F16
  ```

## Running the Example

```bash
cd examples/22-structured-json-output
go run main.go
```

> Ref: https://platform.openai.com/docs/guides/structured-outputs