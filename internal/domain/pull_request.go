package domain

import (
	"errors"
	"time"
)

var ErrPRIsExists = errors.New("PR id already exists")
var ErrPRNotFound = errors.New("PR not found")

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
	MergedAt          *time.Time
}
