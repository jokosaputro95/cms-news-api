package dto

import "time"

type RegisterUserInput struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
	Email    string `json:"email" validate:"required,email,min=3,max=30"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type RegisterUserOutput struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}