-- name: InsertUser :exec
INSERT INTO users (username, user_id) VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE set username = EXCLUDED.username where users.username <> EXCLUDED.username;

-- name: GetUserById :one
SELECT * FROM users WHERE user_id = $1;

-- name: GetUserByName :one
SELECT * FROM users WHERE username = $1;