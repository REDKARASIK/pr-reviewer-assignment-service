DROP INDEX IF EXISTS prs.idx_pr_reviewers_user_id;
DROP INDEX IF EXISTS prs.idx_pr_author;

DROP TABLE IF EXISTS prs.pr_reviewers;
DROP TABLE IF EXISTS prs.pull_requests;

DROP TYPE IF EXISTS prs.pr_status;

