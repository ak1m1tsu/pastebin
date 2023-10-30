package repo

import (
	"context"
	"fmt"

	"github.com/romankravchuk/pastebin/internal/entity"
	"github.com/romankravchuk/pastebin/internal/usecase"
	"github.com/romankravchuk/pastebin/pkg/postgres"
)

var _ usecase.UsersRepo = &UsersRepo{}

type UsersRepo struct {
	*postgres.Postgres
}

func NewUsersRepositry(pg *postgres.Postgres) *UsersRepo {
	return &UsersRepo{pg}
}

// GetByEmail godoc.
func (r *UsersRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	sql, args, err := r.Builder.
		Select("id", "email", "username", "avatar", "access_token").
		From("users").
		Where("email = ?", email).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("UsersRepo.GetByEmail.Builder: %w", err)
	}

	user := &entity.User{}

	err = r.Pool.
		QueryRow(ctx, sql, args...).
		Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.Avatar,
			&user.AccessToken,
		)
	if err != nil {
		return nil, fmt.Errorf("UsersRepo.GetByEmail.Pool: %w", err)
	}

	return user, nil
}

// Create godoc.
func (r *UsersRepo) Create(ctx context.Context, u *entity.User) error {
	sql, args, err := r.Builder.
		Insert("users").
		Columns("username", "email", "avatar", "access_token").
		Values(u.Username, u.Email, u.Avatar, u.AccessToken).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return fmt.Errorf("UsersRepo.Create.Builder: %w", err)
	}

	err = r.Pool.
		QueryRow(ctx, sql, args...).
		Scan(&u.ID)
	if err != nil {
		return fmt.Errorf("UsersRepo.Create.Pool: %w", err)
	}

	return nil
}
