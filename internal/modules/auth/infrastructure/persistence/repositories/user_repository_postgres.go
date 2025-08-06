package repositories

import (
	"context"
	"database/sql"
	"time"

	entities "github.com/jokosaputro95/cms-news-api/internal/modules/auth/domain/entities"
	repos "github.com/jokosaputro95/cms-news-api/internal/modules/auth/domain/repositories"
	vo "github.com/jokosaputro95/cms-news-api/internal/modules/auth/domain/value_objects"
)

type UserRepositoryPostgres struct {
	db *sql.DB
}

func NewUserRepositoryPostgres(db *sql.DB) repos.UserRepository {
	return &UserRepositoryPostgres{db: db}
}

func (r *UserRepositoryPostgres) Save(ctx context.Context, user *entities.User) (*entities.User, error) {
	query := `
		INSERT INTO users (id, username, email, hashed_password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`
	var createdAt, updatedAt time.Time
	
	err := r.db.QueryRowContext(
		ctx,
		query,
		user.ID,
		user.Username.String(),
		user.Email.String(),
		user.HashedPassword,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&createdAt, &updatedAt)
	
	if err != nil {
		return nil, err
	}
	
	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt
	
	return user, nil
}

func (r *UserRepositoryPostgres) Update(ctx context.Context, user *entities.User) (*entities.User, error) {
	query := `
		UPDATE users
		SET username = $2, email = $3, hashed_password = $4, updated_at = $5
		WHERE id = $1
		RETURNING updated_at
	`
	var updatedAt time.Time
	
	err := r.db.QueryRowContext(
		ctx,
		query,
		user.ID,
		user.Username.String(),
		user.Email.String(),
		user.HashedPassword,
		time.Now(),
	).Scan(&updatedAt)
	
	if err != nil {
		return nil, err
	}
	
	user.UpdatedAt = updatedAt
	return user, nil
}

func (r *UserRepositoryPostgres) FindByID(ctx context.Context, id string) (*entities.User, error) {
	query := "SELECT id, username, email, hashed_password, created_at, updated_at FROM users WHERE id = $1"
	
	var user entities.User
	var username, email, password string
	var createdAt, updatedAt time.Time
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&username,
		&email,
		&password,
		&createdAt,
		&updatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil // User tidak ditemukan
	}
	if err != nil {
		return nil, err
	}

	// ✅ Recreate value objects dari data yang diambil
	usernameVO, err := vo.NewUsername(username)
	if err != nil {
		return nil, err
	}
	user.Username = *usernameVO

	emailVO, err := vo.NewEmail(email)
	if err != nil {
		return nil, err
	}
	user.Email = *emailVO

	// ✅ Hashed password langsung sebagai string
	user.HashedPassword = password
	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt
	
	return &user, nil
}

func (r *UserRepositoryPostgres) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := "SELECT id, username, email, hashed_password, created_at, updated_at FROM users WHERE email = $1"
	
	var user entities.User
	var username, emailStr, password string // ✅ Fix: scan email ke string dulu
	var createdAt, updatedAt time.Time
	
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&username,
		&emailStr,
		&password,
		&createdAt,
		&updatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// ✅ Recreate value objects dari data yang diambil
	usernameVO, err := vo.NewUsername(username)
	if err != nil {
		return nil, err
	}
	user.Username = *usernameVO

	emailVO, err := vo.NewEmail(emailStr)
	if err != nil {
		return nil, err
	}
	user.Email = *emailVO

	// ✅ Hashed password langsung sebagai string
	user.HashedPassword = password
	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt
	
	return &user, nil
}

func (r *UserRepositoryPostgres) FindAll(ctx context.Context) ([]*entities.User, error) {
	query := "SELECT id, username, email, hashed_password, created_at, updated_at FROM users"
	
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var users []*entities.User
	for rows.Next() {
		var user entities.User
		var username, email, password string
		var createdAt, updatedAt time.Time
		
		if err := rows.Scan(
			&user.ID,
			&username,
			&email,
			&password,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		// ✅ Recreate value objects
		usernameVO, err := vo.NewUsername(username)
		if err != nil {
			return nil, err
		}
		user.Username = *usernameVO

		emailVO, err := vo.NewEmail(email)
		if err != nil {
			return nil, err
		}
		user.Email = *emailVO

		// ✅ Hashed password sebagai string
		user.HashedPassword = password
		user.CreatedAt = createdAt
		user.UpdatedAt = updatedAt
		
		users = append(users, &user)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return users, nil
}

func (r *UserRepositoryPostgres) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)"
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	
	return exists, nil
}

func (r *UserRepositoryPostgres) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}