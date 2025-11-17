package postgres

import (
	"context"
	"fmt"
	"pr-reviewer-assigment-service/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StatisticsPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewStatisticsPostgresRepository(pool *pgxpool.Pool) *StatisticsPostgresRepository {
	return &StatisticsPostgresRepository{
		pool: pool,
	}
}

// GetUserAssignmentStats возвращает статистику по пользователям с лимитом/оффсетом.
func (r *StatisticsPostgresRepository) GetUserAssignmentStats(
	ctx context.Context,
	limit, offset int,
) (*domain.UserAssignmentStatsPage, error) {
	if limit <= 0 {
		return &domain.UserAssignmentStatsPage{
			Items:  []domain.UserAssignmentStat{},
			Total:  0,
			Limit:  limit,
			Offset: offset,
		}, nil
	}

	const query = `
		WITH stats AS (
			SELECT
				a.user_id,
				u.name,
				COUNT(*) AS assignments_count
			FROM prs.pr_reviewers AS a
			LEFT JOIN users.users AS u ON u.id = a.user_id
			GROUP BY a.user_id, u.name
		)
		SELECT
			stats.user_id,
			stats.name,
			stats.assignments_count,
			COUNT(*) OVER() AS total_count
		FROM stats
		ORDER BY stats.assignments_count DESC
		LIMIT $1 OFFSET $2;
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query user stats: %w", err)
	}
	defer rows.Close()

	stats := make([]domain.UserAssignmentStat, 0, limit)
	total := 0
	firstRow := true

	for rows.Next() {
		var s domain.UserAssignmentStat
		var totalCount int

		if err := rows.Scan(&s.UserID, &s.Username, &s.AssignmentsCount, &totalCount); err != nil {
			return nil, fmt.Errorf("scan user stats: %w", err)
		}

		if firstRow {
			total = totalCount
			firstRow = false
		}

		stats = append(stats, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	if firstRow {
		total = 0
	}

	page := &domain.UserAssignmentStatsPage{
		Items:  stats,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}

	return page, nil
}
