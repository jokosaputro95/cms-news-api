package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	dto "github.com/jokosaputro95/cms-news-api/internal/modules/auth/application/dto"
	usecases "github.com/jokosaputro95/cms-news-api/internal/modules/auth/application/usecases"
)

type AuthHandler struct {
	registerUseCase *usecases.RegisterUser
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func NewAuthHandler(registerUseCase *usecases.RegisterUser) *AuthHandler {
	return &AuthHandler{
		registerUseCase: registerUseCase,
	}
}

// Register handles POST /api/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		h.writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set response content type
	w.Header().Set("Content-Type", "application/json")

	// Parse request body
	var input dto.RegisterUserInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if err := h.validateInput(&input); err != nil {
		h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Execute use case
	result, err := h.registerUseCase.Execute(r.Context(), &input)
	if err != nil {
		statusCode := h.getStatusCodeFromError(err)
		h.writeErrorResponse(w, err.Error(), statusCode)
		return
	}

	// Success response
	h.writeSuccessResponse(w, result, http.StatusCreated)
}

// validateInput basic input validation
func (h *AuthHandler) validateInput(input *dto.RegisterUserInput) error {
	if strings.TrimSpace(input.Username) == "" {
		return &ValidationError{Message: "username is required"}
	}
	if strings.TrimSpace(input.Email) == "" {
		return &ValidationError{Message: "email is required"}
	}
	if strings.TrimSpace(input.Password) == "" {
		return &ValidationError{Message: "password is required"}
	}
	return nil
}

// getStatusCodeFromError determines HTTP status code based on error
func (h *AuthHandler) getStatusCodeFromError(err error) int {
	switch {
	case strings.Contains(err.Error(), "already exists"):
		return http.StatusConflict
	case strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "cannot be empty"):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// writeSuccessResponse writes success response
func (h *AuthHandler) writeSuccessResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	response := SuccessResponse{
		Success: true,
		Message: "Request processed successfully",
		Data:    data,
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// writeErrorResponse writes error response
func (h *AuthHandler) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	response := ErrorResponse{
		Success: false,
		Message: "Request failed",
		Error:   message,
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// ValidationError custom error type
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}