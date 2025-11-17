package statistics

type UserStatItemResponse struct {
	UserID           string `json:"user_id"`
	Username         string `json:"username"`
	AssignmentsCount int    `json:"assignments_count"`
}

type UserStatsResponse struct {
	Items  []UserStatItemResponse `json:"items"`
	Total  int                    `json:"total"`
	Limit  int                    `json:"limit"`
	Offset int                    `json:"offset"`
}
