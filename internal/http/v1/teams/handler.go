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

// Add godoc
// @Summary Создать команду с участниками (создаёт/обновляет пользователей)
// @Description
//   - Если команды ещё нет — создаётся команда и все участники добавляются в team_members.
//   - Если команда уже есть — обновляются участники (добавляются/удаляются) и флаг is_active у пользователей.
//   - Если пользователь уже состоит в другой команде — вернётся ошибка.
//
// @Tags Teams
// @Accept json
// @Produce json
// @Param request body TeamAddRequest true "Команда и её участники"
// @Success 201 {object} TeamAddResponse "Созданная/обновлённая команда"
// @Failure 400 {object} response.ErrorResponse "INVALID_JSON"
// @Failure 404 {object} response.ErrorResponse "NOT_FOUND"
// @Failure 409 {object} response.ErrorResponse "USERS_TEAM_EXISTS"
// @Failure 500 {object} response.ErrorResponse "INTERNAL_ERROR"
// @Router /team/add [post]
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
			response.Error(w, http.StatusBadRequest, "USERS_TEAM_EXISTS", err.Error())
		case errors.Is(err, domain.ErrTeamNotFound):
			response.Error(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		case errors.Is(err, domain.ErrUserNotFound):
			response.Error(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		case errors.Is(err, domain.ErrUserAlreadyInTeam):
			response.Error(w, http.StatusConflict, "TEAMS_CONFLICT", err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}
		return
	}

	var teamResponse TeamAddResponse
	teamResponse.Team.TeamName = team.TeamName
	for _, member := range team.Members {
		teamResponse.Team.Members = append(teamResponse.Team.Members, Member{
			Username: member.Username,
			UserID:   member.UserID,
			IsActive: member.IsActive,
		})
	}

	response.JSON(w, http.StatusCreated, teamResponse)
}

// Get godoc
// @Summary Получить команду с участниками
// @Description Возвращает состав команды по её имени.
// @Tags Teams
// @Accept json
// @Produce json
// @Param team_name query string true "Уникальное имя команды"
// @Success 200 {object} TeamResponse "Команда и её участники"
// @Failure 404 {object} response.ErrorResponse "NOT_FOUND / TEAM_NOT_FOUND"
// @Failure 500 {object} response.ErrorResponse "INTERNAL_ERROR"
// @Router /team/get [get]
func (handler *TeamsHandler) Get(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	queryParams := r.URL.Query()
	teamName := queryParams.Get("team_name")
	if teamName == "" {
		response.Error(w, http.StatusBadRequest, "MISSING_FIELD", "team_name field is required")
		return
	}

	teamDomain, err := handler.teamService.GetTeam(r.Context(), teamName)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTeamNotFound):
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "team not found")
		default:
			response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}
		return
	}

	var teamResponse TeamResponse
	teamResponse.TeamName = teamName

	for _, member := range teamDomain.Members {
		teamResponse.Members = append(teamResponse.Members, Member{
			Username: member.Username,
			UserID:   member.UserID,
			IsActive: member.IsActive,
		})
	}

	response.JSON(w, http.StatusOK, teamResponse)
}
