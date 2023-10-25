package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/romankravchuk/pastebin/internal/entity"
	"github.com/romankravchuk/pastebin/internal/usecase"
	"github.com/romankravchuk/pastebin/pkg/postgres"
)

var _ usecase.PastesRepo = &PastesRepo{}

type PastesRepo struct {
	pg *postgres.Postgres
}

func NewPastesRepo(pg *postgres.Postgres) *PastesRepo {
	return &PastesRepo{pg: pg}
}

// DeletePaste implements usecase.PastesRepo.
func (r *PastesRepo) Delete(ctx context.Context, hash string) error {
	sql, args, err := r.pg.Builder.
		Delete("pastes").
		Where("hash = $1", hash).
		ToSql()
	if err != nil {
		return fmt.Errorf("PastesRepo.DeletePaste.Builder: %w", err)
	}

	_, err = r.pg.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("PastesRepo.DeletePaste.Pool.Exec: %w", err)
	}

	return nil
}

// GetPaste implements usecase.PastesRepo.
func (r *PastesRepo) Get(ctx context.Context, hash string) (*entity.Paste, error) {
	var (
		columns = []string{
			"hash",
			"user_id",
			"title",
			"format",
			"password_hash",
			"expires_at",
			"created_at",
		}
		query = r.pg.Builder.
			Select(columns...).
			From("pastes").
			Where("hash = ?", hash)
	)

	s, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("PastesRepo.GetPaste.Builder: %w", err)
	}

	paste := entity.Paste{}

	err = r.pg.Pool.QueryRow(ctx, s, args...).
		Scan(
			&paste.Hash,
			&paste.UserID,
			&paste.Title,
			&paste.Format,
			&paste.Password.Hash,
			&paste.ExpiresAt,
			&paste.CreatedAt,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, usecase.ErrRecordNotFound
		}

		return nil, fmt.Errorf("PastesRepo.GetPaste.Pool.QueryRow: %w", err)
	}

	return &paste, nil
}

// Create inserts a paste metadata in database and upload paste text in blob storage.
func (r *PastesRepo) Create(ctx context.Context, p *entity.Paste) error {
	var (
		columns = []string{"hash", "format"}
		values  = []any{p.Hash, p.Format}
		query   = r.pg.Builder.Insert("pastes")
	)

	if p.Title != "" {
		columns = append(columns, "title")
		values = append(values, p.Title)
	}

	if p.UserID.Valid {
		columns = append(columns, "user_id")
		values = append(values, p.UserID)
	}

	if p.Password.Hash != nil {
		columns = append(columns, "password_hash")
		values = append(values, p.Password.Hash)
	}

	if !p.ExpiresAt.IsZero() {
		columns = append(columns, "expires_at")
		values = append(values, p.ExpiresAt)
	}

	sql, args, err := query.
		Columns(columns...).
		Values(values...).
		Suffix("RETURNING created_at, expires_at").
		ToSql()
	if err != nil {
		return fmt.Errorf("PastesRepo.CreatePaste.Builder: %w", err)
	}

	err = r.pg.Pool.
		QueryRow(ctx, sql, args...).
		Scan(&p.CreatedAt, &p.ExpiresAt)
	if err != nil {
		return fmt.Errorf("PastesRepo.CreatePaste.Pool.Begin: %w", err)
	}

	return nil
}

// UpdatePaste implements usecase.PastesRepo.
func (*PastesRepo) Update(context.Context, *entity.Paste) error {
	panic("unimplemented")
}
