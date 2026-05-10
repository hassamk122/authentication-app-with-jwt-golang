package errs

import "errors"

var (
	ErrEmailTaken            = errors.New("Email already taken")
	ErrInvalidRequestPayload = errors.New("Invalid request payload")
	ErrInvalidCredentails    = errors.New("Invalid credentials")
)
