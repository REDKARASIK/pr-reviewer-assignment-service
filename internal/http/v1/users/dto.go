package users

type SetActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type UserResponse struct {
	UserID   string  `json:"user_id"`
	Username string  `json:"username"`
	TeamName *string `json:"team_name,omitempty"`
	IsActive bool    `json:"is_active"`
}

type SetIsActiveResponse struct {
	User UserResponse `json:"user"`
}

type PullRequestResponse struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

type GetReviewResponse struct {
	UserID       string                `json:"user_id"`
	PullRequests []PullRequestResponse `json:"pull_requests"`
}
