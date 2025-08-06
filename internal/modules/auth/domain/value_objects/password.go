package valueobjects

import (
	"errors"
	"regexp"
	"strings"
)

type Password struct {
	value string
}

const (
	MinPasswordLength = 8
	MaxPasswordLength = 128
)

var (
	ErrPasswordEmpty             = errors.New("password cannot be empty")
	ErrPasswordInvalidLength     = errors.New("password must be between 8 and 128 characters")
	ErrPasswordInvalidCharacters = errors.New("password cannot contain spaces")
)

var passwordRegex = regexp.MustCompile(`^[^\s]+$`)

func NewPassword(value string) (*Password, error) {
	payload := strings.TrimSpace(value)

	if payload == "" {
		return nil, ErrPasswordEmpty
	}

	if len(payload) < MinPasswordLength || len(payload) > MaxPasswordLength {
		return nil, ErrPasswordInvalidLength
	}

	if !passwordRegex.MatchString(payload) {
		return nil, ErrPasswordInvalidCharacters
	}

	return &Password{value: payload}, nil
}

func (p *Password) String() string {
	return p.value
}

func (p *Password) Value() string {
	return p.value
}