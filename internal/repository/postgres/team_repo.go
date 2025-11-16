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
				return err
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
