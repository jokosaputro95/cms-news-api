package entities_test

import (
	"errors"
	"testing"

	entities "github.com/jokosaputro95/cms-news-api/internal/modules/auth/domain/entities"
	vo "github.com/jokosaputro95/cms-news-api/internal/modules/auth/domain/value_objects"
)

// Helper function
func CreateValidUser(t *testing.T) *entities.User {
	t.Helper()

	usernameVO, err := vo.NewUsername("abc")
	if err != nil {
		t.Fatalf("Error creating username value object: %v", err)
	}

	emailVO, err := vo.NewEmail("abc@test.com")
	if err != nil {
		t.Fatalf("Error creating email value object: %v", err)
	}

	// ✅ Simulasi hashed password (biasanya hasil dari bcrypt)
	hashedPassword := "$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewgBMnVJD.U7p0u."

	payload := struct {
		ID             string
		Username       vo.Username
		Email          vo.Email
		HashedPassword string // ✅ Ubah ke HashedPassword (string)
	}{
		ID:             "user-uuid",
		Username:       *usernameVO,
		Email:          *emailVO,
		HashedPassword: hashedPassword,
	}

	user, err := entities.NewUser(payload.ID, payload.Username, payload.Email, payload.HashedPassword)
	
	if err != nil {
		t.Fatalf("Error creating user: %v", err)
	}

	return user
}

func TestNewUserRegister(t *testing.T) {
	t.Run("should create a new user with valid data", func(t *testing.T) {
		user := CreateValidUser(t)
		expectedHashedPassword := "$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewgBMnVJD.U7p0u."

		if user.ID != "user-uuid" {
			t.Errorf("Expected user ID to be 'user-uuid', but got '%s'", user.ID)
		}

		if user.Username.String() != "abc" {
			t.Errorf("Expected username to be 'abc', but got '%s'", user.Username.String())
		}

		if user.Email.String() != "abc@test.com" {
			t.Errorf("Expected email to be 'abc@test.com', but got '%s'", user.Email.String())
		}

		// ✅ Test hashed password sebagai string
		if user.HashedPassword != expectedHashedPassword {
			t.Errorf("Expected hashed password to be '%s', but got '%s'", expectedHashedPassword, user.HashedPassword)
		}

		if user.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set, but it's zero")
		}

		if user.UpdatedAt.IsZero() {
			t.Error("Expected UpdatedAt to be set, but it's zero")
		}
	})

	// ✅ Test validasi raw password (tetap menggunakan value object untuk validasi)
	t.Run("should return an error for empty username", func(t *testing.T) {
		_, err := vo.NewUsername("")
		if err == nil { t.Error("Expected an error for empty username, but got nil") }
		if !errors.Is(err, vo.ErrUsernameEmpty) { t.Errorf("Expected error ErrUsernameEmpty, but got %v", err) }
	})

	t.Run("should return an error for empty email", func(t *testing.T) {
		_, err := vo.NewEmail("")
		if err == nil { t.Error("Expected an error for empty email, but got nil") }
		if !errors.Is(err, vo.ErrEmailEmpty) { t.Errorf("Expected error ErrEmailEmpty, but got %v", err) }
	})
	
	t.Run("should return an error for empty password", func(t *testing.T) {
		_, err := vo.NewPassword("")
		if err == nil { t.Error("Expected an error for empty password, but got nil") }
		if !errors.Is(err, vo.ErrPasswordEmpty) { t.Errorf("Expected error ErrPasswordEmpty, but got %v", err) }
	})

	t.Run("should return an error for too short password", func(t *testing.T) {
		_, err := vo.NewPassword("1234567") // 7 karakter, kurang dari minimum 8
		if err == nil { t.Error("Expected an error for short password, but got nil") }
		if !errors.Is(err, vo.ErrPasswordInvalidLength) { t.Errorf("Expected error ErrPasswordInvalidLength, but got %v", err) }
	})
}