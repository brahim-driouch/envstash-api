package services

import "errors"

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrWeakPassword       = errors.New("password too weak")
	ErrUserNotFound       = errors.New("no account found with this emails")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnexpected         = errors.New("could not proceed you request")
	ErrUserNotVerified    = errors.New("please verify your email")
)
