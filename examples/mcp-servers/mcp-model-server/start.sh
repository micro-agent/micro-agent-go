#!/bin/bash
export MODEL_RUNNER_BASE_URL=http://localhost:12434/engines/llama.cpp/v1
export CHAT_MODEL=ai/qwen2.5:0.5B-F16
export MCP_HTTP_PORT=9090
export AGENT_NAME="Bob"
export SYSTEM_INSTRUCTION="You are Bob, a helpful AI assistant."
go run main.go