package service

import (
	"context"
	"pr-reviewer-assigment-service/internal/domain"
	"sort"
)

type PullRequestRepository interface {
	GetReviewPRs(ctx context.Context, userID string) ([]domain.PullRequest, error)
	Create(ctx context.Context, prID, prName, authorID string) error
	AssignReviewers(ctx context.Context, prID string, reviewers []string) error
}

type PullRequestService struct {
	repo     PullRequestRepository
	userRepo UserRepository
	teamRepo TeamRepository
}

func NewPullRequestService(repo PullRequestRepository, userRepo UserRepository, teamRepo TeamRepository) *PullRequestService {
	return &PullRequestService{repo: repo, userRepo: userRepo, teamRepo: teamRepo}
}

func (service *PullRequestService) GetReview(ctx context.Context, id string) ([]domain.PullRequest, error) {
	prs, err := service.repo.GetReviewPRs(ctx, id)
	if err != nil {
		return nil, err
	}
	return prs, nil
}

func (service *PullRequestService) Create(ctx context.Context, prID, prName, authorID string) (*domain.PullRequestAssignment, error) {
	user, err := service.userRepo.GetByID(ctx, authorID)
	if err != nil {
		return nil, err
	}
	if user.TeamName == nil {
		return nil, domain.ErrTeamNotFound
	}

	err = service.repo.Create(ctx, prID, prName, authorID)
	if err != nil {
		return nil, err
	}

	assignments, err := service.teamRepo.GetTeamsMembersByTeamName(ctx, *user.TeamName)
	if err != nil {
		return nil, err
	}

	sort.Slice(assignments, func(i, j int) bool {
		var left, right int64
		if assignments[i].PRReviews != nil {
			left = *assignments[i].PRReviews
		}
		if assignments[j].PRReviews != nil {
			right = *assignments[j].PRReviews
		}
		return left < right
	})

	var prAssignments domain.PullRequestAssignment
	for _, assignment := range assignments {
		if len(prAssignments.AssignedReviewers) == 2 {
			break
		}

		if assignment.UserID != authorID {
			prAssignments.AssignedReviewers = append(prAssignments.AssignedReviewers, assignment.UserID)
		}
	}

	err = service.repo.AssignReviewers(ctx, prID, prAssignments.AssignedReviewers)
	if err != nil {
		return nil, err
	}

	prAssignments.PullRequestID = prID
	prAssignments.PullRequestName = prName
	prAssignments.AuthorID = authorID
	prAssignments.Status = domain.PROpenStatus

	return &prAssignments, nil
}
