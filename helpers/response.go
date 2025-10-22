package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

type Response struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

type ErrorResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Error     string `json:"error,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

func SendResponse(w http.ResponseWriter, r *http.Request, data interface{}, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := Response{
		Success:   code >= 200 && code < 300,
		Message:   message,
		Data:      data,
		RequestID: middleware.GetReqID(r.Context()),
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		if Logger != nil {
			Logger.Errorf("Failed to encode response: %v", err)
		}
	}
}

func SendErrorResponse(w http.ResponseWriter, r *http.Request, message string, err error, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
		// Log the error if logger is available
		if Logger != nil {
			Logger.WithFields(map[string]interface{}{
				"error":      err,
				"message":    message,
				"request_id": middleware.GetReqID(r.Context()),
			}).Error("Error response")
		}
	}

	resp := ErrorResponse{
		Success:   false,
		Message:   message,
		Error:     errorMsg,
		RequestID: middleware.GetReqID(r.Context()),
	}

	if encodeErr := json.NewEncoder(w).Encode(resp); encodeErr != nil {
		if Logger != nil {
			Logger.Errorf("Failed to encode error response: %v", encodeErr)
		}
	}
}
