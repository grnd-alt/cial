
-- name: CreateComment :one
INSERT INTO COMMENTS(
    id, post_id, user_id, content, user_name
) VALUES(
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetCommentsByPost :many
select * from comments where post_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3 ;


-- name: GetCommentsByPosts :many
SELECT c.* from unnest($1::varchar[]) as post_ids
JOIN LATERAL (
    SELECT * FROM comments WHERE post_id = post_ids ORDER BY created_at DESC LIMIT 3
) c ON true;

-- name: DeleteCommentsByPost :exec
DELETE FROM COMMENTS WHERE post_id = $1;
