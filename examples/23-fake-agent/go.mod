module fake-agent

go 1.24.0

toolchain go1.24.4

replace github.com/micro-agent/micro-agent-go => ../../

require github.com/micro-agent/micro-agent-go v0.0.0-00010101000000-000000000000

require (
	github.com/openai/openai-go/v2 v2.1.1
	github.com/tidwall/gjson v1.14.4 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/sjson v1.2.5 // indirect
)
