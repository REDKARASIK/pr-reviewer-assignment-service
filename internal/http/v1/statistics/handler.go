package statistics

import (
	"net/http"
	"pr-reviewer-assigment-service/internal/http/response"
	"pr-reviewer-assigment-service/internal/service"
	"strconv"
)

type StatisticsHandler struct {
	statsService *service.StatisticsService
}

func NewStatisticsHandler(statsService *service.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{
		statsService: statsService,
	}
}

// GetUserStats возвращает статистику назначений по пользователям.
// @Summary      Статистика назначений по пользователям
// @Description  Возвращает список пользователей и количество назначенных им PR с пагинацией.
// @Tags         statistics
// @Accept       json
// @Produce      json
// @Param        limit   query     int  false  "Лимит выборки (по умолчанию 50, максимум 100)"
// @Param        offset  query     int  false  "Смещение выборки (по умолчанию 0)"
// @Success      200     {object}  UserStatsResponse
// @Failure      500     {object}  response.ErrorResponse
// @Router       /stats/users [get]
func (h *StatisticsHandler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query()

	limit := 50
	offset := 0

	if limitStr := query.Get("limit"); limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil && v > 0 {
			if v > 100 {
				v = 100
			}
			limit = v
		}
	}

	if offsetStr := query.Get("offset"); offsetStr != "" {
		if v, err := strconv.Atoi(offsetStr); err == nil && v >= 0 {
			offset = v
		}
	}

	ctx := r.Context()

	page, err := h.statsService.GetUserAssignmentStats(ctx, limit, offset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	resp := UserStatsResponse{
		Items:  make([]UserStatItemResponse, 0, len(page.Items)),
		Total:  page.Total,
		Limit:  page.Limit,
		Offset: page.Offset,
	}

	for _, s := range page.Items {
		resp.Items = append(resp.Items, UserStatItemResponse{
			UserID:           s.UserID,
			Username:         s.Username,
			AssignmentsCount: s.AssignmentsCount,
		})
	}

	response.JSON(w, http.StatusOK, resp)
}
