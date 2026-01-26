-- name: CreateCounter :one
INSERT INTO counters (name, icon) VALUES ( $1, $2 ) RETURNING *;

-- name: AddUserToCounter :exec
INSERT INTO counters_users (user_id, counter_id, token, access_type) VALUES ( $1, $2, $3, $4);

-- name: AddEventToCounter :exec
INSERT INTO counters_users_events (user_id, counter_id) VALUES ($1, $2);

-- name: GetCounterEventsForUser :many
SELECT * FROM counters_users_events where user_id = $1 limit $2;

-- name: GetCountersForUser :many
SELECT DISTINCT
  counters.*,
  counters_users.access_type,
  counters_users.entry_count
FROM counters,counters_users
WHERE counters_users.counter_id = counters.id
 AND counters_users.user_id = $1;

-- name: GetUserInCounter :one
SELECT DISTINCT * FROM counters_users where user_id = $1 and counter_id = $2;

-- name: GetUsersInCounter :many
SELECT DISTINCT counters_users.access_type, counters_users.entry_count, users.username from counters_users LEFT JOIN users on counters_users.user_id = users.user_id
where counters_users.counter_id = $1;

-- name: GetCounter :one
SELECT 
    counters.*
FROM 
    counters
WHERE 
    counters.id = $1;

-- name: GetEvents :many
SELECT * from counters_users_events where user_id =$1 and counter_id = $2;
