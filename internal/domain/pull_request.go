package domain

import (
	"errors"
	"time"
)

// ErrPRIsExists возвращается, если PR с таким ID уже существует.
var ErrPRIsExists = errors.New("PR id already exists")

// ErrPRNotFound возвращается, если PR не существует
var ErrPRNotFound = errors.New("PR not found")

// ErrPRMerged возвращается, если PR смержен
var ErrPRMerged = errors.New("cannot reassign on merged PR")

// ErrIsNotAssigned возвращается, если user не назначен для этого PR
var ErrIsNotAssigned = errors.New("reviewer is not assigned to this PR")

// ErrIsNoCandidates возвращается, если нет доступных кандидатов на назначения для PR
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
