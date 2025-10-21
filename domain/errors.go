package domain

import "errors"

var (
	ErrInvalidChatID  = errors.New("invalid chat ID")
	ErrInvalidRequest = errors.New("invalid request")
	ErrNotFound       = errors.New("record not found")
	ErrInvalidEntity  = errors.New("invalid entity")
)
