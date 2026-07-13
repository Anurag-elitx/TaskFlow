package repository

import (
	"context"
	"errors"
	"taskflow/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
}

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (email, password_hash, full_name)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, user.Email, user.PasswordHash, user.FullName).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, email, password_hash, full_name, created_at, updated_at FROM users WHERE email = $1`
	var user models.User
	err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	query := `SELECT id, email, password_hash, full_name, created_at, updated_at FROM users WHERE id = $1`
	var user models.User
	err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Update(ctx context.Context, user *models.User) error {
	query := `UPDATE users SET full_name = $1, updated_at = NOW() WHERE id = $2 RETURNING updated_at`
	return r.db.QueryRow(ctx, query, user.FullName, user.ID).Scan(&user.UpdatedAt)
}
