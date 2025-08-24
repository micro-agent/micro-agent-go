package mu

import "fmt"

// ExitToolCallsLoopError signals early termination of tool call processing loops
type ExitToolCallsLoopError struct {
	Message string
}

// Error implements the error interface for ExitToolCallsLoopError
func (e *ExitToolCallsLoopError) Error() string {
	return fmt.Sprintf("Message: %s", e.Message)
}

// ExitStreamCompletionError signals early termination of streaming completions
type ExitStreamCompletionError struct {
	Message string
}

// Error implements the error interface for ExitStreamCompletionError
func (e *ExitStreamCompletionError) Error() string {
	return fmt.Sprintf("Message: %s", e.Message)
}
