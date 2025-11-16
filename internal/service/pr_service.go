package service

import (
	"context"
	"pr-reviewer-assigment-service/internal/domain"
)

type PullRequestRepository interface {
	GetReviewPRs(ctx context.Context, userID string) ([]domain.PullRequest, error)
}

type PullRequestService struct {
	repo PullRequestRepository
}

func NewPullRequestService(repo PullRequestRepository) *PullRequestService {
	return &PullRequestService{repo: repo}
}

func (service *PullRequestService) GetReview(ctx context.Context, id string) ([]domain.PullRequest, error) {
	prs, err := service.repo.GetReviewPRs(ctx, id)
	if err != nil {
		return nil, err
	}
	return prs, nil
}
