-- +goose Up
-- +goose StatementBegin

CREATE TABLE about (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    bio VARCHAR(255) NOT NULL DEFAULT '',
    company VARCHAR(255) NOT NULL DEFAULT '',
    readme TEXT NOT NULL DEFAULT '',
    job_title VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_about_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT unique_about_user_id
        UNIQUE (user_id)
);

CREATE TABLE social_links (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    github VARCHAR(255) NOT NULL DEFAULT '',
    linkedin VARCHAR(255) NOT NULL DEFAULT '',
    website VARCHAR(255) NOT NULL DEFAULT '',
    x VARCHAR(255) NOT NULL DEFAULT '',
    youtube VARCHAR(255) NOT NULL DEFAULT '',
    stack_overflow VARCHAR(255) NOT NULL DEFAULT '',
    reddit VARCHAR(255) NOT NULL DEFAULT '',
    roadmap_sh VARCHAR(255) NOT NULL DEFAULT '',
    codepen VARCHAR(255) NOT NULL DEFAULT '',
    mastodon VARCHAR(255) NOT NULL DEFAULT '',
    threads VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_social_links_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT unique_social_links_user_id
        UNIQUE (user_id)
);

CREATE TABLE user_stats (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    followers INTEGER NOT NULL DEFAULT 0,
    following INTEGER NOT NULL DEFAULT 0,
    reputation INTEGER NOT NULL DEFAULT 0,
    views INTEGER NOT NULL DEFAULT 0,
    upvotes INTEGER NOT NULL DEFAULT 0,

    CONSTRAINT fk_user_stats_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT unique_user_stats_user_id
        UNIQUE (user_id)
);

CREATE TABLE streaks (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    streak_start TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    streak_end TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    streak_length INTEGER NOT NULL DEFAULT 1,
    longest_streak INTEGER NOT NULL DEFAULT 1,

    CONSTRAINT fk_streaks_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT unique_streaks_user_id
        UNIQUE (user_id)
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS streaks;
DROP TABLE IF EXISTS user_stats;
DROP TABLE IF EXISTS social_links;
DROP TABLE IF EXISTS about;

-- +goose StatementEnd