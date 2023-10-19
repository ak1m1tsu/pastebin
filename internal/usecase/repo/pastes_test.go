package repo

import (
	"testing"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/romankravchuk/pastebin/internal/entity"
	"github.com/romankravchuk/pastebin/pkg/postgres"
	"github.com/stretchr/testify/require"
)

func Test_1(t *testing.T) {
	pg := &postgres.Postgres{Builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar)}
	p := entity.Paste{
		Hash:      "HuO10P",
		Name:      "user-id",
		Format:    "json",
		ExpiresAt: time.Time{},
	}
	p.Password.Set("hello")

	columns := []string{"hash", "name", "format"}
	values := []any{p.Hash, p.Name, p.Format}

	if p.UserID != "" {
		columns = append(columns, "user_id")
		values = append(values, p.UserID)
	}

	if p.Password.Hash() != nil {
		columns = append(columns, "password_hash")
		values = append(values, p.Password.Hash())
	}

	if !p.ExpiresAt.IsZero() {
		columns = append(columns, "expires_at")
		values = append(values, p.ExpiresAt)
	}

	query := pg.Builder.Insert("pastes").
		Columns(columns...).
		Values(values...)

	sql, args, err := query.ToSql()

	require.NoError(t, err)
	t.Log(sql)
	t.Log(args...)
}
