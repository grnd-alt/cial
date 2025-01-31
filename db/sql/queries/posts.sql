-- name: GetOne :one
select * from posts where posts.id = $1 ;

-- name: CreatePost :one
INSERT INTO posts(
    created_by, content, username, id, filepath
) VALUES(
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetPostsByUser :many
SELECT * FROM posts WHERE created_by = $1 OR username = $2 ORDER BY created_at DESC LIMIT $3 OFFSET $4;

-- name: GetLatestPosts :many
SELECT * FROM posts ORDER BY created_at DESC LIMIT $1 OFFSET $2;


-- name: DeletePost :exec
DELETE FROM posts WHERE id = $1;
