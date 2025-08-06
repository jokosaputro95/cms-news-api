package valueobjects

import (
	"errors"
	"regexp"
	"strings"
)

type Email struct {
	value string
}

const (
	MinEmailLength = 3
	MaxEmailLength = 100
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

var (
	ErrEmailEmpty        = errors.New("email cannot be empty")
	ErrEmailInvalid      = errors.New("email is not valid")
	ErrEmailInvalidLength = errors.New("email must be between 3 and 100 characters")
)

func NewEmail(value string) (*Email, error) {
	payload := strings.TrimSpace(value)

	if payload == "" {
		return nil, ErrEmailEmpty
	}

	if len(payload) < MinEmailLength || len(payload) > MaxEmailLength {
		return nil, ErrEmailInvalidLength
	}

	if !emailRegex.MatchString(payload) {
		return nil, ErrEmailInvalid
	}


	return &Email{value: payload}, nil
}

func (e *Email) String() string {
	return e.value
}

func (e *Email) Value() string {
	return e.value
}