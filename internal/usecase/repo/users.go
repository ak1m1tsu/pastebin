package repo

import (
	"context"
	"fmt"

	"github.com/romankravchuk/pastebin/internal/entity"
	"github.com/romankravchuk/pastebin/pkg/postgres"
)

type UsersRepo struct {
	*postgres.Postgres
}

func NewUsersRepo(pg *postgres.Postgres) *UsersRepo {
	return &UsersRepo{pg}
}

func (r *UsersRepo) Create(_ context.Context, u *entity.User) error {
	_, _, err := r.Builder.
		Insert("users").
		Columns("username", "email", "avatar", "access_token").
		Values(u.Username, u.Email, u.Avatar, u.AccessToken).
		Suffix("RETURNING id, created_at").
		ToSql()
	if err != nil {
		return fmt.Errorf("UsersRepo.Create.Builder: %w", err)
	}

	return nil
}
