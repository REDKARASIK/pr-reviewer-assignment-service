package domain

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
