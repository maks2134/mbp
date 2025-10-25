package errors_constant

import "errors"

var (
	UserAlreadyExists        = errors.New("user already exists")
	UserNotFound             = errors.New("user not found")
	UserAuthenticationFailed = errors.New("invalid username or password")
)
