package domain

import (
	"errors"
	"time"
)

var ErrPRIsExists = errors.New("PR id already exists")
var ErrPRNotFound = errors.New("PR not found")
var ErrPRMerged = errors.New("cannot reassign on merged PR")
var ErrIsNotAssigned = errors.New("reviewer is not assigned to this PR")
var ErrIsNoCandidates = errors.New("no active replacement candidate in team")

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
	ReplacedBy        *string
}
