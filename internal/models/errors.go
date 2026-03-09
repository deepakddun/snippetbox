package models

import "errors"

var ErrorNoRecord = errors.New("models: no matching record found")

var ErrInvalidCredentials = errors.New("models: invalid credentails")

var ErrDuplicateEmail = errors.New("models: duplicate email")
