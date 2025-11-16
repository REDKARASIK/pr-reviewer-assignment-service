CREATE SCHEMA IF NOT EXISTS users;

CREATE TABLE IF NOT EXISTS users.users (
    id BIGSERIAL PRIMARY KEY,
    name varchar(125) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users.teams (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users.team_members (
    team_id BIGINT NOT NULL
        REFERENCES users.teams(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL
        REFERENCES users.users(id) ON DELETE CASCADE,
    joined_at timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (team_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_team_members_user_id
    ON users.team_members(user_id);