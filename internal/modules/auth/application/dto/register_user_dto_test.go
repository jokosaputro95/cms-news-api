package dto_test

import (
	"testing"
	"time"

	dto "github.com/jokosaputro95/cms-news-api/internal/modules/auth/application/dto"
)

func TestRegisterUserInput(t *testing.T) {
	t.Run("should be a valid struct", func(t *testing.T) {
		input := dto.RegisterUserInput{
			Username: "testuser",
			Email: "test@example.com",
			Password: "testpassword123",
		}

		if input.Username != "testuser" || input.Email != "test@example.com" || input.Password != "testpassword123" {
			t.Errorf("RegisterUserInput struct is not valid. Expected: %+v, Got: %+v", input, input)
		}
	})
}

func TestRegisterUserOutput(t *testing.T) {
	t.Run("should be a valid struct", func(t *testing.T) {
		output := dto.RegisterUserOutput{
			ID: "123",
			Username: "testuser",
			Email: "test@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		
		if output.CreatedAt.IsZero() || output.UpdatedAt.IsZero() {
			t.Errorf("RegisterUserOutput struct is not valid. Expected: %+v, Got: %+v", output, output)
		}
	})	
}