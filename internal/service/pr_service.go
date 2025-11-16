package service

import (
	"context"
	"pr-reviewer-assigment-service/internal/domain"
	"slices"
	"sort"
)

type PullRequestRepository interface {
	GetReviewPRs(ctx context.Context, userID string) ([]domain.PullRequest, error)
	Create(ctx context.Context, prID, prName, authorID string) error
	AssignReviewers(ctx context.Context, prID string, reviewers []string) error
	Merge(ctx context.Context, prID string) (*domain.PullRequestAssignment, error)
	IsExists(ctx context.Context, prID string) (bool, error)
	GetPRReviewers(ctx context.Context, prID string) ([]string, error)
	GetPRAuthors(ctx context.Context, prID string) (string, error)
	GetPRNameByID(ctx context.Context, prID string) (string, error)
	DeleteAssignedUser(ctx context.Context, prID string) error
}

const PRReviewers int = 2

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
		if len(prAssignments.AssignedReviewers) == PRReviewers {
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

func (service *PullRequestService) Merge(ctx context.Context, prID string) (*domain.PullRequestAssignment, error) {
	prAssignments, err := service.repo.Merge(ctx, prID)
	if err != nil {
		return nil, err
	}

	return prAssignments, nil
}

func (service *PullRequestService) Reassign(ctx context.Context, prID, replacedUserID string) (*domain.PullRequestAssignment, error) {
	isPrExists, err := service.repo.IsExists(ctx, prID)
	if err != nil {
		return nil, err
	}
	if !isPrExists {
		return nil, domain.ErrPRNotFound
	}

	user, err := service.userRepo.GetByID(ctx, replacedUserID)
	if err != nil {
		return nil, err
	}
	prs, err := service.repo.GetReviewPRs(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if hasPR := slices.ContainsFunc(prs, func(pr domain.PullRequest) bool {
		return pr.PullRequestID == prID
	}); !hasPR {
		return nil, domain.ErrIsNotAssigned
	}

	if isPrMerges := slices.ContainsFunc(prs, func(pr domain.PullRequest) bool {
		return pr.PullRequestID == prID && pr.Status == domain.PRMergeStatus
	}); isPrMerges {
		return nil, domain.ErrPRMerged
	}

	reviewers, err := service.repo.GetPRReviewers(ctx, prID)
	if err != nil {
		return nil, err
	}

	authorID, err := service.repo.GetPRAuthors(ctx, prID)
	if err != nil {
		return nil, err
	}
	author, err := service.userRepo.GetByID(ctx, authorID)
	if err != nil {
		return nil, err
	}

	if author.TeamName == nil {
		return nil, domain.ErrTeamNotFound
	}

	assignments, err := service.teamRepo.GetTeamsMembersByTeamName(ctx, *author.TeamName)
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
	existing := make(map[string]struct{})

	for _, reviewer := range reviewers {
		if reviewer != replacedUserID {
			prAssignments.AssignedReviewers = append(prAssignments.AssignedReviewers, reviewer)
			existing[reviewer] = struct{}{}
		}

	}
	countReviews := len(prAssignments.AssignedReviewers)

	for _, assignment := range assignments {
		if len(prAssignments.AssignedReviewers) == PRReviewers {
			break
		}

		if assignment.UserID == replacedUserID || assignment.UserID == authorID {
			continue
		}
		if _, used := existing[assignment.UserID]; used {
			continue
		}

		prAssignments.AssignedReviewers = append(prAssignments.AssignedReviewers, assignment.UserID)
		existing[assignment.UserID] = struct{}{}
	}

	if len(prAssignments.AssignedReviewers) == countReviews {
		return nil, domain.ErrIsNoCandidates
	}

	err = service.repo.DeleteAssignedUser(ctx, prID)
	if err != nil {
		return nil, err
	}

	err = service.repo.AssignReviewers(ctx, prID, prAssignments.AssignedReviewers)
	if err != nil {
		return nil, err
	}

	prName, err := service.repo.GetPRNameByID(ctx, prID)
	if err != nil {
		return nil, err
	}

	prAssignments.PullRequestID = prID
	prAssignments.PullRequestName = prName
	prAssignments.AuthorID = authorID
	prAssignments.Status = domain.PROpenStatus
	prAssignments.ReplacedBy = &replacedUserID

	return &prAssignments, nil
}
