package entity

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrDuplicateEmail = errors.New("duplicate email")

type User struct {
	ID       string `db:"id"`
	Email    string `db:"email"`
	Activate bool   `db:"activated"`
	Password
}

type Password struct {
	Plaintext *string `db:"-"`
	Hash      []byte  `db:"password_hash"`
}

func (p *Password) Set(plaintext string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.Hash = hash
	p.Plaintext = &plaintext

	return nil
}

func (p *Password) Matches(plaintext string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintext))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
