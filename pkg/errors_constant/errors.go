package errors_constant

import "errors"

var (
	UserAlreadyExists = errors.New("user already exists")
	UserNotFound      = errors.New("user not found")
	PostNotFound      = errors.New("post not found")
	PostInvalidInput  = errors.New("invalid input")
)
