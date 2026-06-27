-- name: FindMagicLinkByToken :one
SELECT
  email
FROM
  magic_links
WHERE
  token = $1;

-- name: VerifyUserByEmail :one
UPDATE users
SET
  is_verified = $1
WHERE
  email = $2
RETURNING *;

-- name: DeleteMagicLinksByEmail :exec
DELETE FROM
  magic_links
WHERE
  email = $1;

-- name: CreateAbout :one
INSERT INTO about (
  user_id,
  bio,
  company,
  job_title
)
VALUES (
  $1,
  '',
  '',
  ''
)
RETURNING *;

-- name: CreateSocialLinks :one
INSERT INTO social_links (
  user_id
)
VALUES (
  $1
)
RETURNING *;

-- name: CreateUserStats :one
INSERT INTO user_stats (
  user_id,
  followers,
  following,
  reputation,
  views,
  upvotes
)
VALUES (
  $1,
  0,
  0,
  0,
  0,
  0
)
RETURNING *;

-- name: CreateStreak :one
INSERT INTO streaks (
  user_id
)
VALUES (
  $1
)
RETURNING *;

-- name: FindVerifiedUserByEmail :one
SELECT
  id,
  name,
  username,
  avatar,
  email,
  user_password,
  profession
FROM
  users
WHERE
  email = $1
  AND is_verified = $2;

-- name: FindUserEmail :one
SELECT
  email
FROM
  users
WHERE
  email = $1;

-- name: CreateMagicLink :one
INSERT INTO magic_links (
  email,
  token
)
VALUES (
  $1,
  $2
)
RETURNING *;

-- name: CreateUser :one
INSERT INTO users (
  name,
  username,
  profession,
  email,
  user_password
)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5
)
RETURNING *;

-- name: FindUserByEmail :one
SELECT
  *
FROM
  users
WHERE
  email = $1;

-- name: UpdateRefreshToken :exec
UPDATE users
SET
  refresh_token = $1
WHERE
  email = $2;

-- name: CreateGithubSocialLink :one
INSERT INTO social_links (
  user_id,
  github
)
VALUES (
  $1,
  $2
)
RETURNING *;

-- name: FindUserByID :one
SELECT
  *
FROM
  users
WHERE
  id = $1;