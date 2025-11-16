package postgres

import (
	"context"
	"errors"
	"pr-reviewer-assigment-service/internal/domain"
	"time"

	"github.com/jackc/pgx/v5"
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

func (repo *PullRequestRepository) Merge(ctx context.Context, prID string) (*domain.PullRequestAssignment, error) {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return nil, err
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

	const qSelectPR = `
		SELECT id, title, author_id, status, merged_at
		FROM prs.pull_requests
		WHERE id = $1
	`

	var (
		id       string
		name     string
		authorID string
		status   domain.PRStatus
		mergedAt *time.Time
	)

	err = tx.QueryRow(ctx, qSelectPR, prID).Scan(&id, &name, &authorID, &status, &mergedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPRNotFound
		}
		return nil, err
	}

	if status != domain.PRMergeStatus {
		const qUpdate = `
		UPDATE prs.pull_requests
		SET status = $2,
		    merged_at = NOW()
		WHERE id = $1
	`

		if _, err = tx.Exec(ctx, qUpdate, prID, domain.PRMergeStatus); err != nil {
			return nil, err
		}
	}

	const qReviewers = `
		SELECT user_id
		FROM prs.pr_reviewers
		WHERE pr_id = $1
		ORDER BY user_id
	`

	rows, err := tx.Query(ctx, qReviewers, prID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var rid string
		if err := rows.Scan(&rid); err != nil {
			return nil, err
		}
		reviewers = append(reviewers, rid)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	if mergedAt == nil {
		err = tx.QueryRow(ctx, `SELECT merged_at FROM prs.pull_requests WHERE id = $1`, prID).Scan(&mergedAt)
		if err != nil {
			return nil, err
		}
	}

	return &domain.PullRequestAssignment{
		PullRequest: domain.PullRequest{
			PullRequestID:   prID,
			PullRequestName: name,
			AuthorID:        authorID,
			Status:          domain.PRMergeStatus,
		},
		AssignedReviewers: reviewers,
		MergedAt:          mergedAt,
	}, nil
}
