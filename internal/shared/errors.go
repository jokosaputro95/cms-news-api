package shared

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidInput = errors.New("invalid input")
	ErrDatabaseError = errors.New("database error")
)

func NewValidationError(message string) error {
	return fmt.Errorf("validation error: %s", message)
}

func NewDatabaseError(err error) error {
	return fmt.Errorf("database error: %w", err)
}

func NewConflictError(message string) error {
	return fmt.Errorf("conflict error: %s", message)
}

// ✅ Helper untuk get error code dari error message
func GetErrorCode(err error) string {
	errMsg := strings.ToLower(err.Error())
	
	switch {
	case strings.Contains(errMsg, "validation"):
		return "VALIDATION_ERROR"
	case strings.Contains(errMsg, "conflict"):
		return "CONFLICT_ERROR"
	case strings.Contains(errMsg, "database"):
		return "DATABASE_ERROR"
	case strings.Contains(errMsg, "not found"):
		return "NOT_FOUND"
	default:
		return "INTERNAL_ERROR"
	}
}

// ✅ Helper untuk get user-friendly message
func GetUserMessage(err error) string {
	errMsg := strings.ToLower(err.Error())
	
	switch {
	case strings.Contains(errMsg, "validation"):
		// Extract message after "validation error: "
		parts := strings.Split(err.Error(), "validation error: ")
		if len(parts) > 1 {
			return parts[1]
		}
		return "Invalid input provided"
	case strings.Contains(errMsg, "conflict"):
		// Extract message after "conflict error: "
		parts := strings.Split(err.Error(), "conflict error: ")
		if len(parts) > 1 {
			return parts[1]
		}
		return "Resource already exists"
	case strings.Contains(errMsg, "database"):
		return "System temporarily unavailable"
	default:
		return "An unexpected error occurred"
	}
}