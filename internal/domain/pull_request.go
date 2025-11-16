package domain

import "errors"

var ErrPRIsExists = errors.New("PR id already exists")

type PRStatus string

const (
	PROpenStatus  PRStatus = "OPEN"
	PRMergeStatus PRStatus = "MERGED"
)

type PullRequest struct {
	PullRequestID   string
	PullRequestName string
	AuthorID        string
	Status          PRStatus
}

type PullRequestAssignment struct {
	PullRequest
	AssignedReviewers []string
}
