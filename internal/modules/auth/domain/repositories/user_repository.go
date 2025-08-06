package repositories

import (
	"context"

	entities "github.com/jokosaputro95/cms-news-api/internal/modules/auth/domain/entities"
)

type UserRepository interface {
	Save(ctx context.Context, user *entities.User) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) (*entities.User, error)
	FindByID(ctx context.Context, id string) (*entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	FindAll(ctx context.Context) ([]*entities.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Delete(ctx context.Context, id string) error
}