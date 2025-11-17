package service

import (
	"context"
	"pr-reviewer-assigment-service/internal/domain"
)

type StatisticsRepository interface {
	GetUserAssignmentStats(ctx context.Context, limit, offset int) (*domain.UserAssignmentStatsPage, error)
}

type StatisticsService struct {
	statsRepo StatisticsRepository
}

func NewStatisticsService(statsRepo StatisticsRepository) *StatisticsService {
	return &StatisticsService{
		statsRepo: statsRepo,
	}
}

func (s *StatisticsService) GetUserAssignmentStats(
	ctx context.Context, limit, offset int) (*domain.UserAssignmentStatsPage, error) {
	return s.statsRepo.GetUserAssignmentStats(ctx, limit, offset)
}
