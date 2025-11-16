CREATE SCHEMA IF NOT EXISTS prs;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'pr_status') THEN
        CREATE TYPE prs.pr_status AS ENUM ('OPEN', 'MERGED');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS prs.pull_requests (
    id VARCHAR(255) PRIMARY KEY,
    title TEXT NOT NULL,
    author_id VARCHAR(255) NOT NULL
        REFERENCES users.users(id) ON DELETE RESTRICT,
    status prs.pr_status NOT NULL DEFAULT 'OPEN',
    created_at timestamptz NOT NULL DEFAULT NOW(),
    merged_at timestamptz
);

CREATE TABLE IF NOT EXISTS prs.pr_reviewers (
    pr_id VARCHAR(255) NOT NULL
        REFERENCES prs.pull_requests(id) ON DELETE CASCADE,
    user_id VARCHAR(255) NOT NULL
        REFERENCES users.users(id) ON DELETE RESTRICT,
    assigned_at timestamptz NOT NULL DEFAULT NOW(),
        PRIMARY KEY (pr_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_pr_author ON prs.pull_requests(author_id);

CREATE INDEX IF NOT EXISTS idx_pr_reviewers_user_id ON prs.pr_reviewers(user_id);