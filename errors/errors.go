package errors

import "errors"

// ErrInvalidLogin is an error that can happen when providing the wrong password or username
var ErrInvalidLogin = errors.New("invalid login or password")
