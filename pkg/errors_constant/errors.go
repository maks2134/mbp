package errors_constant

import "errors"

var (
	UserAlreadyExists = errors.New("User already exists")
	UserNotFound      = errors.New("User not found")
)
