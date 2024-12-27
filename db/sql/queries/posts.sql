-- name: GetOne :one
select * from posts where id = $1;

-- name: CreatePost :one
INSERT INTO posts(
    created_by, content, username, id, filepath
) VALUES(
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetPostsByUser :many
SELECT * from posts where created_by = $1 order by created_at desc;
