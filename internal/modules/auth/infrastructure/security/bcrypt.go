package security

import (
	"golang.org/x/crypto/bcrypt"

	vo "github.com/jokosaputro95/cms-news-api/internal/modules/auth/domain/value_objects"
)

type BcryptHasher struct{
	cost int
}

func NewBcryptHasher(cost int) vo.Hasher {
	return &BcryptHasher{cost: cost}
}

func (h *BcryptHasher) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (h *BcryptHasher) Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}