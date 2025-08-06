package entities

import (
	"time"

	vo "github.com/jokosaputro95/cms-news-api/internal/modules/auth/domain/value_objects"
)

type User struct {
	ID string
	Username vo.Username
	Email vo.Email
	HashedPassword string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(id string, username vo.Username, email vo.Email, hashedPassword string) (*User, error) {
	now := time.Now()
	return &User{
		ID: id,
		Username: username,
		Email: email,
		HashedPassword: hashedPassword,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

