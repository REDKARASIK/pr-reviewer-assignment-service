package domain

type UserAssignmentStat struct {
	UserID           string
	Username         string
	AssignmentsCount int
}

type UserAssignmentStatsPage struct {
	Items  []UserAssignmentStat
	Total  int
	Limit  int
	Offset int
}
