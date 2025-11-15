CREATE TYPE pr_status AS ENUM ('OPEN', 'MERGED');

CREATE TABLE IF NOT EXISTS Teams(
    name VARCHAR PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS Users(
    id VARCHAR PRIMARY KEY,
    name TEXT,
    team_name VARCHAR REFERENCES Teams(name),
    isActive BOOLEAN
);

CREATE TABLE IF NOT EXISTS PullRequests(
    id VARCHAR PRIMARY KEY,
    title TEXT,
    author_id VARCHAR REFERENCES Users(id),
    status pr_status NOT NULL DEFAULT 'OPEN',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    merged_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS PullRequestsUsers(
    id VARCHAR PRIMARY KEY,
    pr_id VARCHAR REFERENCES PullRequests(id),
    user_id VARCHAR REFERENCES Users(id)
);
