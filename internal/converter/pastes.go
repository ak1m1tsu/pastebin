package converter

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"time"

	"github.com/romankravchuk/pastebin/internal/entity"
)

const hashLen = 8

func generateHash(text string) string {
	var (
		b       = make([]byte, hashLen)
		r       = rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec // use for randomness only, not security
		hash    = sha256.Sum256([]byte(text))
		encHash = base64.StdEncoding.EncodeToString(hash[:])
	)

	for i := range b {
		b[i] = encHash[r.Intn(len(encHash))]
	}

	return string(b)
}

func CreatePasteToEntity(body *entity.CreatePasteBody) *entity.Paste {
	p := &entity.Paste{
		Hash:      generateHash(body.Text),
		Name:      body.Name,
		Format:    body.Format,
		ExpiresAt: time.Now().Add(2 * 365 * 24 * time.Hour),
		File:      entity.File(body.Text),
	}
	p.Password.Set(body.Password)

	if body.Expires != "" {
		p.ExpiresAt, _ = time.Parse(time.RFC3339, body.Expires) //nolint:errcheck // should always provide valid time
	}

	return p
}

func ModelToResponse(model *entity.Paste) *entity.PasteResponse {
	return &entity.PasteResponse{
		Hash:      model.Hash,
		Title:     model.Name,
		Format:    model.Format,
		ExpiresAt: model.ExpiresAt.Format(time.RFC1123),
		CreatedAt: model.CreatedAt.Format(time.RFC1123),
	}
}
