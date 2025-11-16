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

// Merge godoc
// @Summary Пометить PR как MERGED (идемпотентная операция)
// @Description
//
//	Завершает pull request и помечает его как MERGED.
//	Если PR уже в статусе MERGED — операция идемпотентна:
//	ничего не изменяется, и возвращаются текущие данные PR.
//	Если PR не существует — возвращается ошибка.
//
// @Tags PullRequests
// @Accept json
// @Produce json
// @Param request body MergePRRequest true "Идентификатор PR для merge"
// @Success 200 {object} MergePRResponse "PR успешно помечен как MERGED"
// @Failure 400 {object} response.ErrorResponse "INVALID_JSON"
// @Failure 404 {object} response.ErrorResponse "NOT_FOUND"
// @Failure 500 {object} response.ErrorResponse "INTERNAL_ERROR"
// @Router /pullRequest/merge [post]
func (handler *PullRequestHandler) Merge(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	var request MergePRRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "invalid request body")
		return
	}

	prMergeInfo, err := handler.prService.Merge(r.Context(), request.PullRequestID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrPRNotFound):
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "resource not found")
		default:
			response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		}
		return
	}

	prMergeResponse := MergePRResponse{
		PullRequest: PullRequestResponse{
			PullRequestID:     prMergeInfo.PullRequestID,
			PullRequestName:   prMergeInfo.PullRequestName,
			AuthorID:          prMergeInfo.AuthorID,
			Status:            string(prMergeInfo.Status),
			AssignedReviewers: prMergeInfo.AssignedReviewers,
			MergedAt:          prMergeInfo.MergedAt,
		},
	}

	response.JSON(w, http.StatusOK, prMergeResponse)
}

// Reassign godoc
// @Summary Переназначить ревьювера на другого из его команды
// @Description
//
//	Заменяет конкретного ревьювера в PR на другого участника той же команды.
//	Новый ревьювер выбирается из активных участников команды с минимальным числом назначенных ревью.
//	Автор PR никогда не попадает в список ревьюверов.
//	Если нет доступного кандидата — возвращается ошибка NO_CANDIDATE.
//
// @Tags PullRequests
// @Accept json
// @Produce json
// @Param request body ReassignPRRequest true "PR и старый ревьювер"
// @Success 200 {object} ReassignPRResponse "Успешное переназначение ревьювера"
// @Failure 400 {object} response.ErrorResponse "INVALID_JSON"
// @Failure 404 {object} response.ErrorResponse "NOT_FOUND"
// @Failure 409 {object} response.ErrorResponse "PR_MERGED / NO_CANDIDATE / NOT_ASSIGNED"
// @Failure 500 {object} response.ErrorResponse "INTERNAL_ERROR"
// @Router /pullRequest/reassign [post]
func (handler *PullRequestHandler) Reassign(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	var request ReassignPRRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "invalid request body")
		return
	}

	prAssgs, err := handler.prService.Reassign(r.Context(), request.PullRequestID, request.OldReviewerID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound) || errors.Is(err, domain.ErrPRNotFound):
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "resource not found")
		case errors.Is(err, domain.ErrPRMerged):
			response.Error(w, http.StatusConflict, "PR_MERGED", err.Error())
		case errors.Is(err, domain.ErrIsNotAssigned):
			response.Error(w, http.StatusConflict, "NOT_ASSIGNED", err.Error())
		case errors.Is(err, domain.ErrIsNoCandidates):
			response.Error(w, http.StatusConflict, "NO_CANDIDATE", err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		}
		return
	}

	prAssgsResponse := ReassignPRResponse{
		PullRequest: PullRequestResponse{
			PullRequestID:     prAssgs.PullRequestID,
			PullRequestName:   prAssgs.PullRequestName,
			AuthorID:          prAssgs.AuthorID,
			Status:            string(prAssgs.Status),
			AssignedReviewers: prAssgs.AssignedReviewers,
			ReplacedBy:        prAssgs.ReplacedBy,
		},
	}

	response.JSON(w, http.StatusOK, prAssgsResponse)
}
