package users

type SetActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type SetActiveResponse struct {
	User struct {
		UserID   string `json:"user_id"`
		Username string `json:"username"`
		TeamName string `json:"team_name"`
		IsActive bool   `json:"is_active"`
	} `json:"user"`
}

type PRStatus string

const (
	PROpenStatus  PRStatus = "OPEN"
	PRMergeStatus PRStatus = "MERGED"
)

type PullRequest struct {
	PullRequestID   string   `json:"pull_request_id"`
	PullRequestName string   `json:"pull_request_name"`
	AuthorID        string   `json:"author_id"`
	Status          PRStatus `json:"status"`
}

type GetReviewResponse struct {
	UserID       string        `json:"user_id"`
	PullRequests []PullRequest `json:"pull_requests"`
}
