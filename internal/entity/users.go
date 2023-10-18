package entity

import (
	"errors"
)

var ErrDuplicateEmail = errors.New("duplicate email")

type User struct {
	ID          string `db:"id"`
	Email       string `db:"email"`
	Username    string `db:"username"`
	Avatar      string `db:"avatar"`
	AccessToken []byte `db:"access_token"`
}
