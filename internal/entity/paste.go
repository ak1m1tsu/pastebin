package entity

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"time"
)

var ErrPasteNotFound = errors.New("paste not found")

type Paste struct {
	Hash      string    `db:"hash"`
	UserID    string    `db:"user_id"`
	Name      string    `db:"name"`
	Format    string    `db:"format"`
	CreatedAt time.Time `db:"created_at"`
	ExpiresAt time.Time `db:"expires_at"`
	Password  password  `db:"password"`
	File      File      `db:"file"`
}

type password struct {
	plaintext string `db:"-"`
	hash      []byte `db:"password"`
}

func (p *password) Set(pwd string) {
	passhash := sha256.Sum256([]byte(pwd))
	p.plaintext = pwd
	p.hash = passhash[:]
}

func (p password) Matches(pwd string) bool {
	return [32]byte(p.hash) == sha256.Sum256([]byte(pwd))
}

func (p password) String() string {
	return p.plaintext
}

func (p password) Hash() []byte {
	return bytes.Clone(p.hash)
}

type File []byte

func (f File) Size() int64 {
	return int64(len(f))
}

// @description Тело запроса для создания пасты.
type CreatePasteBody struct {
	// Текст пасты
	Text string `json:"text" example:"this is my paste" validate:"required"`
	// Формат текста
	Format string `json:"format" example:"json" enums:"json,yaml,toml" validate:"required,oneof=json plaintext toml yaml xml"`
	// Время, через которое паста становится не доступной
	Expires string `json:"expires" example:"30m" validate:"omitempty,oneof=30m 1h 1w 1mth"`
	// Пароль для получения доступа к пасте
	Password string `json:"password" example:"hello" validate:"omitempty,max=255"`
	// Название
	Name string `json:"name" example:"thie is paste" validate:"omitempty,max=255"`
} // @name Paste

type PasteResponse struct {
	Hash      string `json:"hash"`
	Title     string `json:"title"`
	Format    string `json:"format"`
	CreatedAt string `json:"created_at"`
	ExpiresAt string `json:"expires_at"`
}
