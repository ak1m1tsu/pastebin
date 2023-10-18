package entity

import (
	"errors"
	"time"
)

var ErrPasteNotFound = errors.New("paste not found")

type Paste struct {
	Hash      string     `db:"hash"`
	UserID    string     `db:"user_id"`
	Name      string     `db:"name"`
	Format    string     `db:"format"`
	URL       string     `db:"url"`
	CreatedAt *time.Time `db:"created_at"`
	ExpiresAt *time.Time `db:"expires_at"`
	Password  []byte     `db:"password"`
	File      File       `db:"file"`
}

type File []byte

func (f File) Size() int64 {
	return int64(len(f))
}
