CREATE TABLE IF NOT EXISTS Teams(
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS Users(
    id SERIAL PRIMARY KEY,
    name TEXT,
    team_id INTEGER REFERENCES Teams(id),
    isActive BOOLEAN
);

CREATE TABLE IF NOT EXISTS PullRequests(
    id SERIAL PRIMARY KEY,
    title TEXT,
    author_id INTEGER REFERENCES Users(id),
    status TEXT,
    needMoreReviewers BOOLEAN
);

CREATE TABLE IF NOT EXISTS PullRequestsUsers(
    id SERIAL PRIMARY KEY,
    pr_id INTEGER REFERENCES PullRequests(id),
    user_id INTEGER REFERENCES Users(id)
);
