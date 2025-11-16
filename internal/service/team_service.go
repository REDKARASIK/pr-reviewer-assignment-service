package service

import (
	"context"
	"pr-reviewer-assigment-service/internal/domain"
)

type TeamRepository interface {
	CreateTeam(ctx context.Context, team *domain.Team) error
	IsTeamExists(ctx context.Context, teamName string) (bool, error)
	UpdateTeamMembers(ctx context.Context, team *domain.Team) error
	GetTeam(ctx context.Context, teamName string) (*domain.Team, error)
	GetTeamsMembersByTeamName(ctx context.Context, teamName string) ([]domain.Member, error)
}

type TeamService struct {
	userRepo UserRepository
	teamRepo TeamRepository
}

func NewTeamService(teamRepo TeamRepository, userRepo UserRepository) *TeamService {
	return &TeamService{
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

func (service *TeamService) Add(ctx context.Context, team domain.Team) (*domain.Team, error) {
	isTeamExists, err := service.teamRepo.IsTeamExists(ctx, team.TeamName)
	if err != nil {
		return nil, err
	}

	var users []domain.User

	for _, member := range team.Members {
		user, err := service.userRepo.GetByID(ctx, member.UserID)
		if err != nil {
			return nil, err
		}

		if isTeamExists && user.TeamName != nil && *user.TeamName != team.TeamName {
			return nil, domain.ErrUserAlreadyInTeam
		}

		if isTeamExists {
			user.IsActive = member.IsActive
			if err = service.userRepo.UpdateActive(ctx, user); err != nil {
				return nil, err
			}
		}

		users = append(users, *user)
	}

	if isTeamExists {
		if err := service.teamRepo.UpdateTeamMembers(ctx, &team); err != nil {
			return nil, err
		}

		updatedTeam := domain.Team{
			TeamName: team.TeamName,
		}
		for _, member := range team.Members {
			updatedTeam.Members = append(updatedTeam.Members, domain.Member{
				UserID:   member.UserID,
				Username: member.Username,
				IsActive: member.IsActive,
			})
		}

		return &updatedTeam, nil
	}

	if err = service.teamRepo.CreateTeam(ctx, &team); err != nil {
		return nil, err
	}

	return &team, nil
}

func (service *TeamService) GetTeam(ctx context.Context, teamName string) (*domain.Team, error) {
	isTeamExists, err := service.teamRepo.IsTeamExists(ctx, teamName)
	if err != nil {
		return nil, err
	}
	if !isTeamExists {
		return nil, domain.ErrTeamNotFound
	}

	team, err := service.teamRepo.GetTeam(ctx, teamName)
	if err != nil {
		return nil, err
	}

	return team, nil
}
