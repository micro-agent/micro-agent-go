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
	"log"
	"net/http"
	"strconv"
)

type A2AServer struct {
	httpPort int
	//Host          string
	httpServer    *http.ServeMux
	agentCard     AgentCard
	agentCallback func(taskRequest TaskRequest) (TaskResponse, error)
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
