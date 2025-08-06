package usecases_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	dto "github.com/jokosaputro95/cms-news-api/internal/modules/auth/application/dto"
	usecases "github.com/jokosaputro95/cms-news-api/internal/modules/auth/application/usecases"
	entities "github.com/jokosaputro95/cms-news-api/internal/modules/auth/domain/entities"
	vo "github.com/jokosaputro95/cms-news-api/internal/modules/auth/domain/value_objects"
)

// --- Mock Implementations ---

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Save(ctx context.Context, user *entities.User) (*entities.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entities.User) (*entities.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*entities.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) FindAll(ctx context.Context) ([]*entities.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.User), args.Error(1)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error{
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockUUIDGenerator struct {
	mock.Mock
}

func (m *MockUUIDGenerator) NewUUID() string {
	args := m.Called()
	return args.String(0)
}

type MockHasher struct {
	mock.Mock
}

func (m *MockHasher) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockHasher) Compare(hashedPassword string, password string) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

// --- Test Suite ---

func setupRegisterUserTest(t *testing.T) (*MockUserRepository, *MockUUIDGenerator, *MockHasher, *usecases.RegisterUser) {
	t.Helper()
	userRepoMock := new(MockUserRepository)
	uuidGenMock := new(MockUUIDGenerator)
	hasherMock := new(MockHasher) 
	
	registerUserUsecase := usecases.NewRegisterUser(userRepoMock, uuidGenMock, hasherMock)
	
	return userRepoMock, uuidGenMock, hasherMock, registerUserUsecase
}

// ✅ Helper function untuk test cases yang tidak perlu semua mock
func setupRegisterUserTestMinimal(t *testing.T) (*MockUserRepository, *MockHasher, *usecases.RegisterUser) {
	t.Helper()
	userRepoMock := new(MockUserRepository)
	uuidGenMock := new(MockUUIDGenerator)
	hasherMock := new(MockHasher)
	
	registerUserUsecase := usecases.NewRegisterUser(userRepoMock, uuidGenMock, hasherMock)
	
	return userRepoMock, hasherMock, registerUserUsecase
}

func TestRegisterUser(t *testing.T) {
	t.Run("should register a user successfully", func(t *testing.T) {
		userRepoMock, uuidGenMock, hasherMock, registerUserUsecase := setupRegisterUserTest(t)
		
		input := dto.RegisterUserInput{
			Username: "jokosaputro",
			Email:    "joko@test.com",
			Password: "password123",
		}
		
		expectedUserID := "mock-uuid-123"
		hashedPassword := "$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewgBMnVJD.U7p0u."

		// Setup mock responses
		userRepoMock.On("ExistsByEmail", mock.Anything, input.Email).Return(false, nil).Once()
		uuidGenMock.On("NewUUID").Return(expectedUserID).Once()
		hasherMock.On("Hash", input.Password).Return(hashedPassword, nil).Once()

		// ✅ Create a mock saved user entity dengan hashed password sebagai string
		usernameVO, _ := vo.NewUsername(input.Username)
		emailVO, _ := vo.NewEmail(input.Email)
		savedUser, _ := entities.NewUser(expectedUserID, *usernameVO, *emailVO, hashedPassword)

		userRepoMock.On("Save", mock.Anything, mock.AnythingOfType("*entities.User")).Return(savedUser, nil).Once()
		
		output, err := registerUserUsecase.Execute(context.Background(), &input)
		
		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, expectedUserID, output.ID)
		assert.Equal(t, input.Username, output.Username)
		assert.Equal(t, input.Email, output.Email)
		assert.NotZero(t, output.CreatedAt) // ✅ Pastikan CreatedAt ada
		assert.NotZero(t, output.UpdatedAt) // ✅ Pastikan UpdatedAt ada
		
		userRepoMock.AssertExpectations(t)
		uuidGenMock.AssertExpectations(t)
		hasherMock.AssertExpectations(t)
	})

	t.Run("should return an error if email already exists", func(t *testing.T) {
		userRepoMock, _, _, registerUserUsecase := setupRegisterUserTest(t)
		
		input := dto.RegisterUserInput{
			Username: "jokosaputro",
			Email:    "joko@test.com",
			Password: "password123",
		}
		
		userRepoMock.On("ExistsByEmail", mock.Anything, input.Email).Return(true, nil).Once()
		
		_, err := registerUserUsecase.Execute(context.Background(), &input)
		
		assert.NotNil(t, err)
		assert.EqualError(t, err, "email already exists")
		
		userRepoMock.AssertExpectations(t)
	})
	
	t.Run("should return an error for invalid input data", func(t *testing.T) {
		_, _, _, registerUserUsecase := setupRegisterUserTest(t)
		
		input := dto.RegisterUserInput{
			Username: "", // Invalid username
			Email:    "joko@test.com",
			Password: "password123",
		}
		
		_, err := registerUserUsecase.Execute(context.Background(), &input)
		
		assert.NotNil(t, err)
		assert.EqualError(t, err, "username cannot be empty")
	})

	// ✅ Test tambahan untuk password validation
	t.Run("should return an error for invalid password", func(t *testing.T) {
		_, _, _, registerUserUsecase := setupRegisterUserTest(t)
		
		input := dto.RegisterUserInput{
			Username: "jokosaputro",
			Email:    "joko@test.com",
			Password: "123", // Too short password
		}
		
		_, err := registerUserUsecase.Execute(context.Background(), &input)
		
		assert.NotNil(t, err)
		assert.EqualError(t, err, "password must be between 8 and 128 characters")
	})

	t.Run("should return an error when hashing fails", func(t *testing.T) {
		userRepoMock, hasherMock, registerUserUsecase := setupRegisterUserTestMinimal(t) // ✅ Gunakan minimal setup
		
		input := dto.RegisterUserInput{
			Username: "jokosaputro",
			Email:    "joko@test.com",
			Password: "password123",
		}
		
		userRepoMock.On("ExistsByEmail", mock.Anything, input.Email).Return(false, nil).Once()
		hasherMock.On("Hash", input.Password).Return("", assert.AnError).Once()
		
		_, err := registerUserUsecase.Execute(context.Background(), &input)
		
		assert.NotNil(t, err)
		assert.Equal(t, assert.AnError, err)
		
		userRepoMock.AssertExpectations(t)
		hasherMock.AssertExpectations(t)
	})

	// ✅ Test case tambahan untuk coverage yang lebih baik
	t.Run("should return an error when repository save fails", func(t *testing.T) {
		userRepoMock, uuidGenMock, hasherMock, registerUserUsecase := setupRegisterUserTest(t)
		
		input := dto.RegisterUserInput{
			Username: "jokosaputro",
			Email:    "joko@test.com",
			Password: "password123",
		}
		
		expectedUserID := "mock-uuid-123"
		hashedPassword := "$2a$12$hashedpassword"

		userRepoMock.On("ExistsByEmail", mock.Anything, input.Email).Return(false, nil).Once()
		uuidGenMock.On("NewUUID").Return(expectedUserID).Once()
		hasherMock.On("Hash", input.Password).Return(hashedPassword, nil).Once()
		userRepoMock.On("Save", mock.Anything, mock.AnythingOfType("*entities.User")).Return(nil, assert.AnError).Once()
		
		_, err := registerUserUsecase.Execute(context.Background(), &input)
		
		assert.NotNil(t, err)
		assert.Equal(t, assert.AnError, err)
		
		userRepoMock.AssertExpectations(t)
		uuidGenMock.AssertExpectations(t)
		hasherMock.AssertExpectations(t)
	})

	t.Run("should return an error when repository ExistsByEmail fails", func(t *testing.T) {
		userRepoMock, _, _, registerUserUsecase := setupRegisterUserTest(t)
		
		input := dto.RegisterUserInput{
			Username: "jokosaputro",
			Email:    "joko@test.com",
			Password: "password123",
		}
		
		userRepoMock.On("ExistsByEmail", mock.Anything, input.Email).Return(false, assert.AnError).Once()
		
		_, err := registerUserUsecase.Execute(context.Background(), &input)
		
		assert.NotNil(t, err)
		assert.Equal(t, assert.AnError, err)
		
		userRepoMock.AssertExpectations(t)
	})
}