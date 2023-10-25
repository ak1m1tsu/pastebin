package entity

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

var ErrPasteNotFound = errors.New("paste not found")

type Paste struct {
	Hash      string         `db:"hash"`
	UserID    sql.NullString `db:"user_id"`
	Title     string         `db:"title"`
	Format    string         `db:"format"`
	CreatedAt time.Time      `db:"created_at"`
	ExpiresAt time.Time      `db:"expires_at"`
	File      File
	Password  Password
}

func (p *Paste) UnmarshalBinary(raw []byte) error {
	return json.Unmarshal(raw, &p)
}

func (p *Paste) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

type Password struct {
	Plaintext string `db:"-"`
	Hash      []byte `db:"password_hash"`
}

func (p *Password) Set(pwd string) {
	if pwd == "" {
		return
	}

	passhash := sha256.Sum256([]byte(pwd))
	p.Plaintext = pwd
	p.Hash = passhash[:]
}

func (p Password) Matches(pwd string) bool {
	return [32]byte(p.Hash) == sha256.Sum256([]byte(pwd))
}

type File []byte

func (f File) Size() int64 {
	return int64(len(f))
}

// @description Тело запроса для создания пасты.
type CreatePasteBody struct {
	// Текст
	Text string `json:"text" example:"Some very secret text" validate:"required"`
	// Формат текста
	Format string `json:"format" example:"plaintext" enums:"json,yaml,toml" validate:"required,oneof=json plaintext toml yaml xml"`
	// Время, через которое паста становится не доступной
	Expires string `json:"expires" example:"30m" validate:"omitempty,oneof=30m 1h 168h 5040h"`
	// Пароль для получения доступа к пасте
	Password string `json:"password" example:"password for security" validate:"omitempty,max=255"`
	// Название
	Title string `json:"title" example:"The private paste" validate:"omitempty,max=255"`
} // @name CreatePasteBody

// @description Тело ответа на создание пасты.
type PasteResponse struct {
	// Уникальный идентификатор
	Hash string `json:"hash"`
	// Название
	Title string `json:"title,omitempty"`
	// Текст
	Text string `json:"text"`
	// Формат текста
	Format string `json:"format"`
	// Дата создания
	CreatedAt string `json:"created_at"`
	// Время, через которое паста становится не доступной
	ExpiresAt string `json:"expires_at"`
} // @name PasteResponse

// @description Тело запроса для разблокировки пасты.
type UnlockPasteBody struct {
	// Пароль
	Password string `json:"password" example:"hello" validate:"required,max=255"`
} // @name UnlockPasteBody
