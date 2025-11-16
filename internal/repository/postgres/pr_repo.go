package postgres

import (
	"context"
	"pr-reviewer-assigment-service/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PullRequestRepository struct {
	pool *pgxpool.Pool
}

func NewPullRequestRepository(pool *pgxpool.Pool) *PullRequestRepository {
	return &PullRequestRepository{
		pool: pool,
	}
}

func (repo *PullRequestRepository) GetReviewPRs(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	const qGetReviews = `
		SELECT
			prr.pr_id,
			pr.title,
			pr.author_id,
			pr.status
		FROM prs.pr_reviewers as prr
		JOIN prs.pull_requests pr ON pr.id = prr.pr_id
		WHERE prr.user_id = $1
	`

	var prs []domain.PullRequest
	rows, err := repo.pool.Query(ctx, qGetReviews, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var pr domain.PullRequest
		err = rows.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status)
		if err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return prs, nil
}
