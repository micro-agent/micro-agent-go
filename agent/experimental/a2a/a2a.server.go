// Package a2a provides experimental functionality for µ-agent.
//
// WARNING: This package is experimental and subject to change.
// The API may change or be removed in future versions without notice.
// Use at your own risk in production environments.
// NOTE: This is a partial implementation of the A2A protocol.
// IMPORTANT: This is a work in progress and may not cover all aspects of the A2A protocol.
package a2a

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type A2AServer struct {
	httpPort int
	//Host          string
	httpServer          *http.ServeMux
	agentCard           AgentCard
	agentCallback       func(taskRequest TaskRequest) (TaskResponse, error)
	agentStreamCallback func(taskRequest TaskRequest, streamFunc func(content string) error) error
}

// NewA2AServer creates a new A2A server with the given parameters
func NewA2AServer(port int, agentCard AgentCard, agentCallback func(taskRequest TaskRequest) (TaskResponse, error)) *A2AServer {
	mux := http.NewServeMux()
	server := &A2AServer{
		httpPort: port,
		//Host:          host,
		httpServer:    mux,
		agentCard:     agentCard,
		agentCallback: agentCallback,
	}
	// Register handlers
	mux.HandleFunc("/.well-known/agentcard", server.getAgentCard)
	mux.HandleFunc("/.well-known/agent.json", server.getAgentCard)
	//mux.HandleFunc("/task", server.handleTaskSync) // Using the synchronous handler
	mux.HandleFunc("/", server.handleTaskSync) // Using the synchronous handler

	return server
}

// NewA2AServerWithStreaming creates a new A2A server with streaming support
func NewA2AServerWithStreaming(port int, agentCard AgentCard, agentCallback func(taskRequest TaskRequest) (TaskResponse, error), agentStreamCallback func(taskRequest TaskRequest, streamFunc func(content string) error) error) *A2AServer {
	mux := http.NewServeMux()
	server := &A2AServer{
		httpPort: port,
		//Host:          host,
		httpServer:          mux,
		agentCard:           agentCard,
		agentCallback:       agentCallback,
		agentStreamCallback: agentStreamCallback,
	}
	// Register handlers
	mux.HandleFunc("/.well-known/agentcard", server.getAgentCard)
	mux.HandleFunc("/.well-known/agent.json", server.getAgentCard)
	//mux.HandleFunc("/task", server.handleTaskSync)   // Using the synchronous handler
	//mux.HandleFunc("/stream", server.handleTaskStream) // Using the streaming handler
	
	//mux.HandleFunc("/", server.handleTaskSync)       // Default to synchronous handler
	mux.HandleFunc("/", server.handleTaskStream)       // Default to synchronous handler

	return server
}

func (a2asvr *A2AServer) Start() error {
	errListening := http.ListenAndServe(":"+strconv.Itoa(a2asvr.httpPort), a2asvr.httpServer)
	if errListening != nil {
		return errListening
	}
	return nil
}

// Serve the Agent Card at the well-known URL
func (a2asvr *A2AServer) getAgentCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(a2asvr.agentCard)
}

// Alternative synchronous implementation that should work better
func (a2asvr *A2AServer) handleTaskSync(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var taskRequest TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&taskRequest); err != nil {
		http.Error(w, `{"error": "invalid request format"}`, http.StatusBadRequest)
		return
	}

	switch taskRequest.Method {
	case "message/send":
		if len(taskRequest.Params.Message.Parts) > 0 {
			// Process the task synchronously without mutex in the HTTP handler
			// The mutex should only be in the AgentCallback if needed
			responseTask, err := a2asvr.agentCallback(taskRequest)
			if err != nil {
				log.Printf("Agent callback failed for task %s: %v", taskRequest.ID, err)
				http.Error(w, `{"error": "agent callback failed"}`, http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(responseTask)
		} else {
			http.Error(w, `{"error": "invalid request format"}`, http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, `{"error": "unknown method"}`, http.StatusBadRequest)
	}
}

// handleTaskStream handles streaming requests using Server-Sent Events (SSE)
func (a2asvr *A2AServer) handleTaskStream(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if a2asvr.agentStreamCallback == nil {
		http.Error(w, `{"error": "streaming not supported"}`, http.StatusMethodNotAllowed)
		return
	}

	var taskRequest TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&taskRequest); err != nil {
		http.Error(w, `{"error": "invalid request format"}`, http.StatusBadRequest)
		return
	}

	switch taskRequest.Method {
	case "message/send":
		if len(taskRequest.Params.Message.Parts) > 0 {
			// Set up Server-Sent Events headers
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")

			// Send initial response with task metadata
			initialResponse := TaskResponse{
				ID:             taskRequest.ID,
				JSONRpcVersion: "2.0",
				Result: Result{
					Status: TaskStatus{
						State: "streaming",
					},
					Kind:     "task",
					Metadata: map[string]any{},
				},
			}

			initialData, _ := json.Marshal(initialResponse)
			fmt.Fprintf(w, "data: %s\n\n", initialData)
			w.(http.Flusher).Flush()

			// Collect streamed content
			var fullContent string

			// Stream callback function
			streamFunc := func(content string) error {
				if content != "" {
					fullContent += content
					// Send streaming chunk
					chunkResponse := map[string]any{
						"id":      taskRequest.ID,
						"type":    "chunk",
						"content": content,
					}
					chunkData, _ := json.Marshal(chunkResponse)
					fmt.Fprintf(w, "data: %s\n\n", chunkData)
					w.(http.Flusher).Flush()
				}
				return nil
			}

			// Call the streaming callback
			err := a2asvr.agentStreamCallback(taskRequest, streamFunc)
			if err != nil {
				log.Printf("Agent stream callback failed for task %s: %v", taskRequest.ID, err)
				errorResponse := map[string]any{
					"id":    taskRequest.ID,
					"type":  "error",
					"error": "agent callback failed",
				}
				errorData, _ := json.Marshal(errorResponse)
				fmt.Fprintf(w, "data: %s\n\n", errorData)
				w.(http.Flusher).Flush()
				return
			}

			// Send final response
			finalResponse := TaskResponse{
				ID:             taskRequest.ID,
				JSONRpcVersion: "2.0",
				Result: Result{
					Status: TaskStatus{
						State: "completed",
					},
					History: []AgentMessage{
						{
							Role: "assistant",
							Parts: []TextPart{
								{
									Text: fullContent,
									Type: "text",
								},
							},
						},
					},
					Kind:     "task",
					Metadata: map[string]any{},
				},
			}

			finalData, _ := json.Marshal(finalResponse)
			fmt.Fprintf(w, "data: %s\n\n", finalData)
			fmt.Fprintf(w, "event: close\ndata: \n\n")
			w.(http.Flusher).Flush()

		} else {
			http.Error(w, `{"error": "invalid request format"}`, http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, `{"error": "unknown method"}`, http.StatusBadRequest)
	}
}
