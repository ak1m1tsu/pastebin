package entity

import (
	"errors"
	"time"
)

var ErrPasteNotFound = errors.New("paste not found")

type Paste struct {
	Hash         string
	UserID       string
	Name         string
	Format       string
	URL          string
	PasswordHash []byte
	CreatedAt    time.Time
	ExpiresAt    time.Time
	Block        []byte
}
