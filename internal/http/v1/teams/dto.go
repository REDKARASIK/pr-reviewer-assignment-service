package teams

type Member struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type TeamAddRequest struct {
	TeamName string   `json:"team_name"`
	Members  []Member `json:"members"`
}

type TeamAddResponse struct {
	Team TeamResponse `json:"team"`
}

type TeamResponse struct {
	TeamName string   `json:"team_name"`
	Members  []Member `json:"members"`
}
