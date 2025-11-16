package teams

import (
	"encoding/json"
	"errors"
	"net/http"
	"pr-reviewer-assigment-service/internal/domain"
	"pr-reviewer-assigment-service/internal/http/response"
	"pr-reviewer-assigment-service/internal/service"
)

type TeamsHandler struct {
	teamService *service.TeamService
}

func NewTeamsHandler(teamService *service.TeamService) *TeamsHandler {
	return &TeamsHandler{teamService: teamService}
}

func (handler *TeamsHandler) Add(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	var request TeamAddRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "invalid request body")
		return
	}

	teamDomain := domain.Team{
		TeamName: request.TeamName,
	}

	for _, member := range request.Members {
		teamDomain.Members = append(teamDomain.Members, domain.Member{
			Username: member.Username,
			UserID:   member.UserID,
			IsActive: member.IsActive,
		})
	}

	team, err := handler.teamService.Add(r.Context(), teamDomain)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTeamAlreadyExists):
			response.Error(w, http.StatusBadRequest, "TEAM_EXISTS", err.Error())
		case errors.Is(err, domain.ErrTeamNotFound):
			response.Error(w, http.StatusBadRequest, "NOT_FOUND", err.Error())
		case errors.Is(err, domain.ErrUserNotFound):
			response.Error(w, http.StatusBadRequest, "NOT_FOUND", err.Error())
		case errors.Is(err, domain.ErrUserNotFoundInTeam):
			response.Error(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}
		return
	}

	var teamResponse TeamAddResponse
	teamResponse.TeamName = team.TeamName
	for _, member := range team.Members {
		teamResponse.Members = append(teamResponse.Members, Member{
			Username: member.Username,
			UserID:   member.UserID,
			IsActive: member.IsActive,
		})
	}

	response.JSON(w, http.StatusCreated, teamResponse)
}

func (handler *TeamsHandler) Get(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

}
