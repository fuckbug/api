package httputils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	v "github.com/go-playground/validator/v10"
)

type ListResponse[T any] struct {
	// Total number of items
	// example: 10
	Count int `json:"count"`

	// Array of items
	Items []T `json:"items"`
}

func NewListResponse[T any](count int, items []T) ListResponse[T] {
	return ListResponse[T]{
		Count: count,
		Items: items,
	}
}

type ErrorResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if payload == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Failed to encode JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func RespondWithPlainError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, ErrorResponse{
		Code:    code,
		Message: message,
		Details: nil,
	})
}

func RespondWithError(w http.ResponseWriter, code int, message string, details map[string]string) {
	RespondWithJSON(w, code, ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	})
}

func HandleValidatorError(w http.ResponseWriter, err error) {
	var ve v.ValidationErrors
	if errors.As(err, &ve) {
		details := make(map[string]string)
		for _, e := range ve {
			details[e.Field()] = e.Tag()
		}
		RespondWithError(w, http.StatusBadRequest, "Validation failed", details)
	}
}
