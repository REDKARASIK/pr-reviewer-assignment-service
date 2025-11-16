package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"pr-reviewer-assigment-service/internal/domain"
	"pr-reviewer-assigment-service/internal/http/response"
	"pr-reviewer-assigment-service/internal/service"
)

type UsersHandler struct {
	userService *service.UserService
	prService   *service.PullRequestService
}

func NewUsersHandler(userService *service.UserService, prService *service.PullRequestService) *UsersHandler {
	return &UsersHandler{
		userService: userService,
		prService:   prService,
	}
}

// SetIsActive
// @Summary      Установить флаг активности пользователя
// @Description  Принимает user_id и is_active, обновляет пользователя и возвращает его состояние
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request  body      SetActiveRequest        true  "Тело запроса"
// @Success      200      {object}  SetIsActiveResponse     "Обновлённый пользователь"
// @Failure      400      {object}  response.ErrorResponse  "Некорректный запрос"
// @Failure      404      {object}  response.ErrorResponse  "Пользователь не найден"
// @Failure      500      {object}  response.ErrorResponse  "Внутренняя ошибка сервера"
// @Router       /users/setIsActive [post]
func (handler *UsersHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	var request SetActiveRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "invalid request body")
		return
	}

	if request.UserID == "" {
		response.Error(w, http.StatusBadRequest, "MISSING_FIELD", "user_id field is required")
		return
	}

	user, err := handler.userService.SetIsActive(r.Context(), request.UserID, request.IsActive)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "resource not found")
		default:
			// общий 500 на всякий случай
			response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}
		return
	}

	resp := SetIsActiveResponse{
		User: UserResponse{
			UserID:   user.ID,
			Username: user.Username,
			TeamName: user.TeamName,
			IsActive: user.IsActive,
		},
	}

	response.JSON(w, http.StatusOK, resp)
}

// GetReview
// @Summary      Получить PR'ы, где пользователь назначен ревьювером
// @Description  Возвращает список PR'ов, в которых user_id указан как ревьювер
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user_id  query     string                 true  "Идентификатор пользователя"
// @Success      200      {object}  GetReviewResponse      "Список PR'ов пользователя"
// @Failure      400      {object}  response.ErrorResponse "Некорректный запрос (нет user_id)"
// @Failure      500      {object}  response.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /users/getReview [get]
func (handler *UsersHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	queryParams := r.URL.Query()
	userID := queryParams.Get("user_id")
	if userID == "" {
		response.Error(w, http.StatusBadRequest, "MISSING_FIELD", "user_id field is required")
		return
	}

	prs, err := handler.prService.GetReview(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	var reviewResponse GetReviewResponse
	reviewResponse.UserID = userID

	for _, pr := range prs {
		prResponse := PullRequestResponse{
			PullRequestID:   pr.PullRequestID,
			PullRequestName: pr.PullRequestName,
			AuthorID:        pr.AuthorID,
			Status:          string(pr.Status),
		}

		reviewResponse.PullRequests = append(reviewResponse.PullRequests, prResponse)
	}

	response.JSON(w, http.StatusOK, reviewResponse)
}
