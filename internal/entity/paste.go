package entity

import (
	"crypto/sha256"
	"encoding/json"
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
	Text string `json:"text" example:"this is my paste" validate:"required"`
	// Формат текста
	Format string `json:"format" example:"json" enums:"json,yaml,toml" validate:"required,oneof=json plaintext toml yaml xml"`
	// Время, через которое паста становится не доступной
	Expires string `json:"expires" example:"30m" validate:"omitempty,oneof=30m 1h 1w 1mth"`
	// Пароль для получения доступа к пасте
	Password string `json:"password" example:"hello" validate:"omitempty,max=255"`
	// Название
	Name string `json:"name" example:"thie is paste" validate:"omitempty,max=255"`
} // @name CreatePasteBody

// @description Тело ответа на создание пасты.
type PasteResponse struct {
	// Уникальный идентификатор
	Hash string `json:"hash"`
	// Название
	Title string `json:"title,omitempty"`
	// Формат текста
	Format string `json:"format"`
	// Дата создания
	CreatedAt string `json:"created_at"`
	// Время, через которое паста становится не доступной
	ExpiresAt string `json:"expires_at"`
} // @name PasteResponse
