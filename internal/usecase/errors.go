package usecase

import "errors"

var (
	ErrPasteNotFound  = errors.New("the paste not found")
	ErrRecordNotFound = errors.New("the record not found")
)
