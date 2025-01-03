
-- name: CreateComment :one
INSERT INTO COMMENTS(
    id, post_id, user_id, content
) VALUES(
    $1, $2, $3, $4
) RETURNING *;

-- name: GetCommentsByPost :many
select * from comments where post_id = $1;

-- name: GetCommentsByPosts :many
WITH RankedComments AS (
    SELECT
        *,
        ROW_NUMBER() OVER (PARTITION BY post_id ORDER BY created_at DESC) AS rn
    FROM comments
    WHERE post_id = ANY($1::varchar[])
)
select * from RankedComments where rn <= 10;
