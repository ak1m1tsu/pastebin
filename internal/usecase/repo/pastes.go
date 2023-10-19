package repo

import (
	"bytes"
	"context"
	"fmt"

	"github.com/romankravchuk/pastebin/internal/entity"
	"github.com/romankravchuk/pastebin/internal/usecase"
	"github.com/romankravchuk/pastebin/pkg/minio"
	"github.com/romankravchuk/pastebin/pkg/postgres"
)

var _ usecase.PastesRepo = &PastesRepo{}

type PastesRepo struct {
	m  *minio.Minio
	pg *postgres.Postgres
}

func NewPastesRepo(pg *postgres.Postgres, m *minio.Minio) *PastesRepo {
	return &PastesRepo{pg: pg, m: m}
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
	sql, args, err := r.pg.Builder.
		Select("hash, user_id, name, format, password_hash, created_at, expires_at").
		From("pastes").
		Where("hash = $1", hash).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("PastesRepo.GetPaste.Builder: %w", err)
	}

	paste := new(entity.Paste)

	err = r.pg.Pool.QueryRow(ctx, sql, args...).
		Scan(
			&paste.Hash,
			&paste.UserID,
			&paste.Name,
			&paste.Format,
			&paste.Password,
			&paste.CreatedAt,
			&paste.ExpiresAt,
		)
	if err != nil {
		return nil, fmt.Errorf("PastesRepo.GetPaste.Pool.QueryRow: %w", err)
	}

	return paste, nil
}

// Create inserts a paste metadata in database and upload paste text in blob storage.
func (r *PastesRepo) Create(ctx context.Context, p *entity.Paste) error {
	var (
		columns = []string{"hash", "format"}
		values  = []any{p.Hash, p.Format}
		query   = r.pg.Builder.Insert("pastes")
	)

	if p.Name != "" {
		columns = append(columns, "name")
		values = append(values, p.Name)
	}

	if p.UserID != "" {
		columns = append(columns, "user_id")
		values = append(values, p.UserID)
	}

	hash := p.Password.Hash()
	if hash != nil {
		columns = append(columns, "password_hash")
		values = append(values, hash)
	}

	if !p.ExpiresAt.IsZero() {
		columns = append(columns, "expires_at")
		values = append(values, p.ExpiresAt)
	}

	sql, args, err := query.
		Columns(columns...).
		Values(values...).
		Suffix("RETURNING created_at").
		ToSql()
	if err != nil {
		return fmt.Errorf("PastesRepo.CreatePaste.Builder: %w", err)
	}

	tx, err := r.pg.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("PastesRepo.CreatePaste.Pool.Begin: %w", err)
	}

	defer tx.Rollback(ctx) //nolint:errcheck // skip errors for rollback - OK

	err = tx.QueryRow(ctx, sql, args...).Scan(&p.CreatedAt)
	if err != nil {
		return fmt.Errorf("PastesRepo.CreatePaste.Pool.QueryRow: %w", err)
	}

	if p.UserID == "" {
		p.UserID = "public"
	}

	err = r.m.UploadObject(ctx, p.UserID, p.Hash, p.File.Size(), bytes.NewReader(p.File))
	if err != nil {
		return fmt.Errorf("PastesRepo.CreatePaste.UploadObject: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("PastesRepo.CreatePaste.Pool.Commit: %w", err)
	}

	return nil
}

// UpdatePaste implements usecase.PastesRepo.
func (*PastesRepo) Update(context.Context, *entity.Paste) error {
	panic("unimplemented")
}
