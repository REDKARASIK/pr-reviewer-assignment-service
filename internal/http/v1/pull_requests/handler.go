package pull_requests

import (
	"encoding/json"
	"errors"
	"net/http"
	"pr-reviewer-assigment-service/internal/domain"
	"pr-reviewer-assigment-service/internal/http/response"
	"pr-reviewer-assigment-service/internal/service"
)

type PullRequestHandler struct {
	prService *service.PullRequestService
}

func NewPullRequestHandler(prService *service.PullRequestService) *PullRequestHandler {
	return &PullRequestHandler{prService: prService}
}

// Create godoc
// @Summary Создать PR и автоматически назначить до 2 ревьюверов из команды автора
// @Description
//
//	Создаёт pull request и выбирает до двух ревьюверов из команды автора с минимальным количеством уже назначенных ревью.
//	Автор PR никогда не попадает в список ревьюверов.
//
// @Tags PullRequests
// @Accept json
// @Produce json
// @Param request body CreatePRRequest true "Параметры для создания PR"
// @Success 201 {object} CreatePRResponse "Созданный PR с назначенными ревьюверами"
// @Failure 400 {object} response.ErrorResponse"INVALID_JSON"
// @Failure 409 {object} response.ErrorResponse "PR_EXISTS"
// @Failure 404 {object} response.ErrorResponse "NOT_FOUND (author or team not found)"
// @Failure 500 {object} response.ErrorResponse "INTERNAL_ERROR"
// @Router /pullRequest/create [post]
func (handler *PullRequestHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	var request CreatePRRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "invalid request body")
		return
	}

	prInfo, err := handler.prService.Create(r.Context(), request.PullRequestID, request.PullRequestName, request.AuthorID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrPRIsExists):
			response.Error(w, http.StatusConflict, "PR_EXISTS", err.Error())
		case errors.Is(err, domain.ErrUserNotFound):
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "resource not found")
		case errors.Is(err, domain.ErrTeamNotFound):
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "resource not found")
		default:
			response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		}
		return
	}

	prResponse := CreatePRResponse{
		PullRequest: PullRequestResponse{
			PullRequestID:     prInfo.PullRequestID,
			PullRequestName:   prInfo.PullRequestName,
			AuthorID:          prInfo.AuthorID,
			Status:            string(prInfo.Status),
			AssignedReviewers: prInfo.AssignedReviewers,
		},
	}

	response.JSON(w, http.StatusCreated, prResponse)
}
