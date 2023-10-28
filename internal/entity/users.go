package entity

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type key int

const (
	UserIDKey key = iota
)

var ErrDuplicateEmail = errors.New("duplicate email")

type User struct {
	ID          string      `db:"id"`
	Email       string      `db:"email"`
	Username    string      `db:"username"`
	Avatar      string      `db:"avatar"`
	AccessToken AccessToken `db:"access_token"`
}

type AccessToken []byte

func (t *AccessToken) GenerateFrom(token []byte) error {
	hash, err := bcrypt.GenerateFromPassword(token, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	*t = hash

	return nil
}

func (t AccessToken) Matches(token []byte) bool {
	return bcrypt.CompareHashAndPassword(t, token) != nil
}

// @description Payload for creating a new user if not exists and get access token.
type CreateTokenRequest struct {
	// Github oauth2 code
	Code string `json:"code"`
} // @name CreateTokenRequest

// @description Payload for getting user info.
type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar_url"`
} // @name UserInfo

// @description Payload for getting access token.
type TokenCredentails struct {
	UserID      string
	Email       string
	AccessToken string
	Type        string
	ExpireAt    time.Time
} // @name TokenCredentials

type APIUser struct {
	Username string `json:"login"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar_url"`
}
