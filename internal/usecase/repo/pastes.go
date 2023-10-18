package repo

import (
	"context"
	"fmt"

	"github.com/romankravchuk/pastebin/internal/entity"
	"github.com/romankravchuk/pastebin/internal/usecase"
	"github.com/romankravchuk/pastebin/pkg/postgres"
)

var _ usecase.PastesRepo = &PastesRepo{}

type PastesRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *PastesRepo {
	return &PastesRepo{pg}
}

// DeletePaste implements usecase.PastesRepo.
func (r *PastesRepo) Delete(ctx context.Context, hash string) error {
	sql, args, err := r.Builder.
		Delete("pastes").
		Where("hash = $1", hash).
		ToSql()
	if err != nil {
		return fmt.Errorf("PastesRepo.DeletePaste.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("PastesRepo.DeletePaste.Pool.Exec: %w", err)
	}

	return nil
}

// GetPaste implements usecase.PastesRepo.
func (r *PastesRepo) Get(ctx context.Context, hash string) (*entity.Paste, error) {
	sql, args, err := r.Builder.
		Select("hash, user_id, name, format, url, password_hash, created_at, expires_at").
		From("pastes").
		Where("hash = $1", hash).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("PastesRepo.GetPaste.Builder: %w", err)
	}

	paste := new(entity.Paste)

	err = r.Pool.QueryRow(ctx, sql, args...).
		Scan(
			&paste.Hash,
			&paste.UserID,
			&paste.Name,
			&paste.Format,
			&paste.URL,
			&paste.Password,
			&paste.CreatedAt,
			&paste.ExpiresAt,
		)
	if err != nil {
		return nil, fmt.Errorf("PastesRepo.GetPaste.Pool.QueryRow: %w", err)
	}

	return paste, nil
}

// Store implements usecase.PastesRepo.
func (r *PastesRepo) Create(ctx context.Context, p *entity.Paste) error {
	sql, args, err := r.Builder.
		Insert("pastes").
		Columns("hash", "user_id", "name", "format", "url", "password_hash", "expires_at").
		Values(p.Hash, p.UserID, p.Name, p.Format, p.URL, p.Password, p.ExpiresAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("PastesRepo.Store.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)

	return err
}

// UpdatePaste implements usecase.PastesRepo.
func (*PastesRepo) Update(context.Context, *entity.Paste) error {
	panic("unimplemented")
}
