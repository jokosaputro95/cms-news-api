package usecases

import (
	"context"

	dto "github.com/jokosaputro95/cms-news-api/internal/modules/auth/application/dto"
	entities "github.com/jokosaputro95/cms-news-api/internal/modules/auth/domain/entities"
	repos "github.com/jokosaputro95/cms-news-api/internal/modules/auth/domain/repositories"
	vo "github.com/jokosaputro95/cms-news-api/internal/modules/auth/domain/value_objects"
	shared "github.com/jokosaputro95/cms-news-api/internal/shared"
)

// RegisterUser adalah use case untuk mendaftarkan pengguna baru.
type RegisterUser struct {
	userRepository repos.UserRepository
	uuidGenerator  shared.UUIDGenerator
	hasher         vo.Hasher
}

// NewRegisterUser adalah konstruktor untuk use case ini.
func NewRegisterUser(
	userRepo repos.UserRepository, 
	uuidGen shared.UUIDGenerator, 
	hasher vo.Hasher) *RegisterUser {
	return &RegisterUser{
		userRepository: userRepo,
		uuidGenerator:  uuidGen,
		hasher:         hasher,
	}
}

// Execute menjalankan logika bisnis pendaftaran pengguna.
func (r *RegisterUser) Execute(ctx context.Context, input *dto.RegisterUserInput) (*dto.RegisterUserOutput, error) {
	// 1. Validasi Input - Validasi raw password terlebih dahulu
	usernameVO, err := vo.NewUsername(input.Username)
	if err != nil {
		return nil, shared.NewValidationError(err.Error())
	}

	emailVO, err := vo.NewEmail(input.Email)
	if err != nil {
		return nil, shared.NewValidationError(err.Error())
	}

	// ✅ Validasi raw password sebelum di-hash
	_, err = vo.NewPassword(input.Password)
	if err != nil {
		return nil, shared.NewValidationError(err.Error())
	}

	// 2. Memeriksa apakah email sudah terdaftar
	isExist, err := r.userRepository.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return nil, shared.NewDatabaseError(err)
	}
	if isExist {
		return nil, shared.NewConflictError("Email already registered")
	}

	// 3. Hash password setelah validasi
	hashedPassword, err := r.hasher.Hash(input.Password)
	if err != nil {
		return nil, shared.NewDatabaseError(err)
	}

	// 4. Membuat Entity User baru dengan hashed password sebagai string
	newUserID := r.uuidGenerator.NewUUID()
	user, err := entities.NewUser(newUserID, *usernameVO, *emailVO, hashedPassword)
	if err != nil {
		return nil, shared.NewValidationError(err.Error())
	}

	// 5. Menyimpan User ke repository
	savedUser, err := r.userRepository.Save(ctx, user)
	if err != nil {
		return nil, shared.NewDatabaseError(err)
	}

	// 6. Mengembalikan DTO output
	output := &dto.RegisterUserOutput{
		ID:        savedUser.ID,
		Username:  savedUser.Username.String(),
		Email:     savedUser.Email.String(),
		CreatedAt: savedUser.CreatedAt,
		UpdatedAt: savedUser.UpdatedAt,
	}

	return output, nil
}