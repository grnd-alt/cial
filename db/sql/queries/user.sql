-- name: InsertUser :exec
INSERT INTO users (username, user_id) VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE set username = EXCLUDED.username where users.username <> EXCLUDED.username;

-- name: GetUserById :one
SELECT * FROM users WHERE user_id = $1;

-- name: GetUserByName :one
SELECT * FROM users WHERE username = $1;

-- name: SetLastLogin :exec
UPDATE users SET last_login = now() WHERE user_id = $1;

-- name: GetNoLoggedInSince :many
SELECT * FROM users WHERE last_login < $1 limit $2;

-- name: SetLastNotified :exec
UPDATE users SET last_notified = now() WHERE user_id = $1;

-- name: FindUser :many
SELECT users.user_id, users.username FROM users WHERE username LIKE '%' || $1 || '%' LIMIT 10;
