package domain

import "errors"

// ErrUserNotFound возвращается, если пользователь с указанным идентификатором отсутствует в системе.
var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID       string
	Username string
	TeamName *string
	IsActive bool
}
