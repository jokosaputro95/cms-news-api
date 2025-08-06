package valueobjects

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

type Username struct {
	value string
}

const (
	MinUsernameLength = 3
	MaxUsernameLength = 30
)

// ^[a-zA-Z] - Harus dimulai dengan huruf.
// [a-zA-Z0-9_]+$ - Diikuti oleh satu atau lebih huruf, angka, atau underscore.
var usernameRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]+$`)

var (
	ErrUsernameEmpty            = errors.New("username cannot be empty")
	ErrUsernameInvalidLength    = errors.New("username must be between 3 and 30 characters")
	ErrUsernameInvalidCharacters  = errors.New("username can only contain letters, numbers, and underscores, and must start with a letter")
)

func NewUsername(value string) (*Username, error) {
	payload := strings.TrimSpace(value)

	if payload == "" {
		return nil, ErrUsernameEmpty
	}

	if len(payload) < MinUsernameLength || len(payload) > MaxUsernameLength {
		return nil, ErrUsernameInvalidLength
	}

	if !usernameRegex.MatchString(payload) {
		return nil, ErrUsernameInvalidCharacters
	}

	firstChar := rune(payload[0])
	if !unicode.IsLetter(firstChar) {
		return nil, ErrUsernameInvalidCharacters
	}

	return &Username{value: payload}, nil
	
}

func (u *Username) String() string {
	return u.value
}

func (u *Username) Value() string {
	return u.value
}