
-- name: GetFollowing :many
SELECT *
FROM user_follows
WHERE follower_id = $1 ORDER BY followed_at DESC LIMIT $2 OFFSET $3;

-- name: GetAllFollowing :many
SELECT *
FROM user_follows
WHERE follower_id = $1 ORDER BY followed_at;

-- name: GetFollowingCount :one
SELECT count(*)
FROM user_follows
WHERE follower_id = $1;


-- name: GetAllFollowers :many
SELECT *
FROM user_follows
WHERE followed_id = $1 ORDER BY followed_at; 

-- name: GetFollowers :many
SELECT *
FROM user_follows
WHERE followed_id = $1 ORDER BY followed_at DESC LIMIT $2 OFFSET $3;

-- name: GetFollowersCount :one
SELECT count(*)
FROM user_follows
WHERE followed_id = $1;


-- name: InsertFollower :exec
INSERT INTO user_follows(follower_id, followed_id, notification_type)
VALUES ($1, $2, $3) ON CONFLICT (follower_id, followed_id) DO NOTHING;

-- name: DeleteFollower :exec
DELETE FROM user_follows
WHERE follower_id = $1 AND followed_id = $2;

-- name: IsFollowing :one
SELECT EXISTS (
    SELECT 1
    FROM user_follows
    WHERE follower_id = $1 AND followed_id = $2
);

-- name: InsertSubscription :exec
INSERT INTO user_subscriptions (user_id, subscription)
VALUES ($1, $2) ON CONFLICT (subscription) DO UPDATE SET user_id=$1;

-- name: DeleteSubscription :exec
DELETE FROM user_subscriptions
WHERE user_id = $1 AND subscription = $2;

-- name: GetSubscriptions :many
SELECT subscription
FROM user_subscriptions
WHERE user_id = $1;