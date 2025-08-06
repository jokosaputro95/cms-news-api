package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	dto "github.com/jokosaputro95/cms-news-api/internal/modules/auth/application/dto"
	usecases "github.com/jokosaputro95/cms-news-api/internal/modules/auth/application/usecases"
	shared "github.com/jokosaputro95/cms-news-api/internal/shared"
)

type AuthHandler struct {
	registerUseCase *usecases.RegisterUser
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewAuthHandler(registerUseCase *usecases.RegisterUser) *AuthHandler {
	return &AuthHandler{
		registerUseCase: registerUseCase,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeErrorResponse(w, "METHOD_NOT_ALLOWED", "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var input dto.RegisterUserInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		h.writeErrorResponse(w, "INVALID_JSON", "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if strings.TrimSpace(input.Username) == "" {
		h.writeErrorResponse(w, "VALIDATION_ERROR", "Username is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(input.Email) == "" {
		h.writeErrorResponse(w, "VALIDATION_ERROR", "Email is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(input.Password) == "" {
		h.writeErrorResponse(w, "VALIDATION_ERROR", "Password is required", http.StatusBadRequest)
		return
	}

	// Execute use case
	result, err := h.registerUseCase.Execute(r.Context(), &input)
	if err != nil {
		errorCode := shared.GetErrorCode(err)
		userMessage := shared.GetUserMessage(err)
		statusCode := h.getStatusCodeFromErrorCode(errorCode)
		
		h.writeErrorResponse(w, errorCode, userMessage, statusCode)
		return
	}

	h.writeSuccessResponse(w, result, "User registered successfully", http.StatusCreated)
}

func (h *AuthHandler) getStatusCodeFromErrorCode(errorCode string) int {
	switch errorCode {
	case "VALIDATION_ERROR":
		return http.StatusBadRequest
	case "CONFLICT_ERROR":
		return http.StatusConflict
	case "DATABASE_ERROR":
		return http.StatusInternalServerError
	case "NOT_FOUND":
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func (h *AuthHandler) writeSuccessResponse(w http.ResponseWriter, data interface{}, message string, statusCode int) {
	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) writeErrorResponse(w http.ResponseWriter, code, message string, statusCode int) {
	response := Response{
		Success: false,
		Message: "Request failed",
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}