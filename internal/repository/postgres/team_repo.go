package postgres

import (
	"context"
	"errors"
	"pr-reviewer-assigment-service/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepository struct {
	pool *pgxpool.Pool
}

func NewTeamRepository(pool *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{pool: pool}
}

func (repo *TeamRepository) CreateTeam(ctx context.Context, team *domain.Team) (err error) {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	const qCreateTeam = `
		INSERT INTO users.teams (name)
		VALUES ($1)
		RETURNING id
	`

	var teamID int64
	if err = tx.QueryRow(ctx, qCreateTeam, team.TeamName).Scan(&teamID); err != nil {
		return err
	}

	const qInsertMember = `
		INSERT INTO users.team_members (team_id, user_id)
		VALUES ($1, $2)
	`

	for _, member := range team.Members {
		if _, err = tx.Exec(ctx, qInsertMember, teamID, member.UserID); err != nil {
			return err
		}
	}

	return nil
}

func (repo *TeamRepository) IsTeamExists(ctx context.Context, teamName string) (bool, error) {
	const qExistsTeam = `
		SELECT 
			EXISTS(
				SELECT
					1
				FROM users.teams
				WHERE name = $1
			);
	`
	var exists bool
	err := repo.pool.QueryRow(ctx, qExistsTeam, teamName).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, domain.ErrTeamNotFound
		}
	}
	return exists, nil
}

// UpdateTeamMembers - обновляет участников команды (UPSERT)
func (repo *TeamRepository) UpdateTeamMembers(ctx context.Context, team *domain.Team) error {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	const qGetTeamID = `
		SELECT id
		FROM users.teams
		WHERE name = $1
	`
	var teamID int64
	if err = tx.QueryRow(ctx, qGetTeamID, team.TeamName).Scan(&teamID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrTeamNotFound
		}
		return err
	}

	const qGetMembers = `
		SELECT user_id
		FROM users.team_members
		WHERE team_id = $1
	`

	rows, err := tx.Query(ctx, qGetMembers, teamID)
	if err != nil {
		return err
	}
	defer rows.Close()

	existing := make(map[string]struct{})
	for rows.Next() {
		var userID string
		if err = rows.Scan(&userID); err != nil {
			return err
		}
		existing[userID] = struct{}{}
	}
	if err = rows.Err(); err != nil {
		return err
	}

	desired := make(map[string]struct{})
	for _, m := range team.Members {
		desired[m.UserID] = struct{}{}
	}

	const qInsertMember = `
		INSERT INTO users.team_members (team_id, user_id)
		VALUES ($1, $2)
	`
	for _, m := range team.Members {
		if _, ok := existing[m.UserID]; !ok {
			if _, err = tx.Exec(ctx, qInsertMember, teamID, m.UserID); err != nil {
				return domain.ErrUserAlreadyInTeam
			}
		}
	}

	const qDeleteMember = `
		DELETE FROM users.team_members
		WHERE team_id = $1 AND user_id = $2
	`
	for userID := range existing {
		if _, ok := desired[userID]; !ok {
			if _, err = tx.Exec(ctx, qDeleteMember, teamID, userID); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetTeam возвращает полную информацию о команде и её участников
func (repo *TeamRepository) GetTeam(ctx context.Context, teamName string) (*domain.Team, error) {
	const qSelectTeam = `
		SELECT id
		FROM users.teams
		WHERE name = $1
	`

	var teamID int64
	err := repo.pool.QueryRow(ctx, qSelectTeam, teamName).Scan(&teamID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrTeamNotFound
		}
		return nil, err
	}

	const qSelectMembers = `
		SELECT u.id, u.name, u.is_active
		FROM users.team_members tm
		JOIN users.users u ON u.id = tm.user_id
		WHERE tm.team_id = $1
		ORDER BY u.name
	`

	rows, err := repo.pool.Query(ctx, qSelectMembers, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []domain.Member
	for rows.Next() {
		var m domain.Member
		if err := rows.Scan(&m.UserID, &m.Username, &m.IsActive); err != nil {
			return nil, err
		}
		members = append(members, m)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &domain.Team{
		TeamName: teamName,
		Members:  members,
	}, nil
}

// GetTeamsMembersByTeamName - возвращает участников команды по её названию
func (repo *TeamRepository) GetTeamsMembersByTeamName(ctx context.Context, teamName string) ([]domain.Member, error) {
	const qSelectTeamID = `
        SELECT id
        FROM users.teams
        WHERE name = $1
    `

	var teamID int64
	err := repo.pool.QueryRow(ctx, qSelectTeamID, teamName).Scan(&teamID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrTeamNotFound
		}
		return nil, err
	}

	const qSelectMembers = `
        SELECT 
            u.id,
            u.name,
            u.is_active,
            (
                SELECT COUNT(*)
                FROM prs.pr_reviewers pra
                WHERE pra.user_id = u.id
            ) AS pr_reviews
        FROM users.team_members tm
        JOIN users.users u ON u.id = tm.user_id
        WHERE tm.team_id = $1 AND u.is_active = true
        ORDER BY u.name;
    `

	rows, err := repo.pool.Query(ctx, qSelectMembers, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []domain.Member

	for rows.Next() {
		var (
			id        string
			username  string
			isActive  bool
			prReviews int64
		)

		if err := rows.Scan(&id, &username, &isActive, &prReviews); err != nil {
			return nil, err
		}

		prPtr := prReviews
		members = append(members, domain.Member{
			UserID:    id,
			Username:  username,
			IsActive:  isActive,
			PRReviews: &prPtr,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}
