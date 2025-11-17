package domain

import "errors"

// ErrTeamNotFound возвращается, если команда с указанным именем не найдена.
var ErrTeamNotFound = errors.New("team not found")

// ErrTeamAlreadyExists возвращается при попытке создать команду,
// которая уже существует в системе.
var ErrTeamAlreadyExists = errors.New("team_name already exists")

// ErrUserAlreadyInTeam возвращается, если пользователь уже состоит в команде.
var ErrUserAlreadyInTeam = errors.New("user already in team")

type Member struct {
	UserID    string
	Username  string
	IsActive  bool
	PRReviews *int64
}

type Team struct {
	TeamName string
	Members  []Member
}
