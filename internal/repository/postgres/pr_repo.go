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

func (repo *PullRequestRepository) Create(ctx context.Context, prID, prName, authorID string) error {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		}
		err = tx.Commit(ctx)
	}()

	const qCreatePR = `INSERT INTO prs.pull_requests(id, title, author_id) VALUES ($1, $2, $3)`

	_, err = tx.Exec(ctx, qCreatePR, prID, prName, authorID)

	if err != nil {
		return domain.ErrPRIsExists
	}

	return nil
}

func (repo *PullRequestRepository) AssignReviewers(ctx context.Context, prID string, reviewers []string) error {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		}
		err = tx.Commit(ctx)
	}()

	const qInsertReviewer = `
		INSERT INTO prs.pr_reviewers (pr_id, user_id)
		VALUES ($1, $2)
	`
	for _, reviewerID := range reviewers {
		if _, err = tx.Exec(ctx, qInsertReviewer, prID, reviewerID); err != nil {
			return err
		}
	}

	return nil
}
