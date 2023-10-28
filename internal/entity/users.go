package entity

import (
	"errors"

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
