package entity

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrDuplicateEmail = errors.New("duplicate email")

type User struct {
	ID          string `db:"id"`
	Email       string `db:"email"`
	Username    string `db:"username"`
	Avatar      string `db:"avatar"`
	accessToken []byte `db:"access_token"`
}

func (u *User) AccessToken() []byte {
	return u.accessToken
}

func (u *User) SetAccessToken(token []byte) error {
	hash, err := bcrypt.GenerateFromPassword(token, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.accessToken = hash

	return nil
}
