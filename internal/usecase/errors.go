package usecase

import "errors"

var (
	ErrPasteNotFound  = errors.New("the paste not found")
	ErrRecordNotFound = errors.New("the record not found")
	ErrNotPasteAuthor = errors.New("the user is not paste authro")
)
