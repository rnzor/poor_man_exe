package cli

import (
	"encoding/json"
	"fmt"

	"github.com/gliderlabs/ssh"
)

// Response is a generic container for API outputs
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// WriteJSON writes a structured response to the SSH session
func WriteJSON(sess ssh.Session, success bool, message string, data interface{}, err error) {
	resp := Response{
		Success: success,
		Message: message,
		Data:    data,
	}
	if err != nil {
		resp.Error = err.Error()
	}

	encoder := json.NewEncoder(sess)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(resp); err != nil {
		fmt.Fprintf(sess, "Error encoding JSON: %v\n", err)
	}
}

// HasFlag checks if a flag exists in the arguments
func HasFlag(args []string, flag string) bool {
	for _, arg := range args {
		if arg == flag {
			return true
		}
	}
	return false
}
