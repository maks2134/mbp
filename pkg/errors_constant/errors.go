package errors_constant

import "errors"

var (
	UserAlreadyExists = errors.New("user already exists")
	UserNotFound      = errors.New("user not found")
	PostNotFound      = errors.New("post not found")
	InvalidTitle      = errors.New("title must be at least 3 characters long")
	UserNotAuthorized = errors.New("user not authorized to modify this post")
)
