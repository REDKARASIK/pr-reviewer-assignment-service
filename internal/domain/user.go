package domain

import "errors"

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID       string
	Username string
	TeamName string
	IsActive bool
}
