package postgres

import (
	"context"
	"errors"
	"pr-reviewer-assigment-service/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (repo *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	const qGetUserById = `
		SELECT 
			u.id as user_id,
			u.name as username,
			u.is_active,
			t.name as team_name
		FROM users.users u
		LEFT JOIN users.team_members tm ON tm.user_id = u.id
		LEFT JOIN users.teams t ON t.id = tm.team_id 
		WHERE u.id = $1
		LIMIT 1
	`

	user := &domain.User{}

	err := repo.pool.QueryRow(ctx, qGetUserById, id).Scan(&user.ID, &user.Username, &user.IsActive, &user.TeamName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (repo *UserRepository) UpdateActive(ctx context.Context, user *domain.User) error {
	const qUpdateActiveByID = `
		UPDATE users.users
		SET is_active = $1
		WHERE id = $2
	`

	cmdTag, err := repo.pool.Exec(ctx, qUpdateActiveByID, user.IsActive, user.ID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (repo *UserRepository) Create(ctx context.Context, userID, name string, isActive *bool) (*domain.User, error) {
	const q = `
		INSERT INTO users.users (id, name, is_active)
		VALUES ($1, $2, $3)
		RETURNING id, name, is_active
	`

	var u domain.User

	err := repo.pool.QueryRow(ctx, q, userID, name, isActive).Scan(
		&u.ID,
		&u.Username,
		&u.IsActive,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
