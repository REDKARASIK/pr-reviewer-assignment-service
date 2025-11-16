package domain

import "errors"

var ErrTeamNotFound = errors.New("team not found")
var ErrTeamAlreadyExists = errors.New("team_name already exists")
var ErrUserNotFoundInTeam = errors.New("user not found in team")

type Member struct {
	UserID   string
	Username string
	IsActive bool
}

type Team struct {
	TeamName string
	Members  []Member
}
