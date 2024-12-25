-- name: GetOne :one
select * from notes where id = $1;

-- name: CreateNote :one
INSERT INTO NOTES(
    created_by, title, content
) VALUES(
    $1, $2, $3
) RETURNING *;

-- name: GetNotesByUser :many
SELECT * from NOTES where created_by = $1 order by created_at desc;
