// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: posts.sql

package dbgen

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createPost = `-- name: CreatePost :one
INSERT INTO posts(
    created_by, content, username, id, filepath
) VALUES(
    $1, $2, $3, $4, $5
) RETURNING id, created_by, username, content, created_at, updated_at, filepath
`

type CreatePostParams struct {
	CreatedBy string
	Content   string
	Username  string
	ID        string
	Filepath  pgtype.Text
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRow(ctx, createPost,
		arg.CreatedBy,
		arg.Content,
		arg.Username,
		arg.ID,
		arg.Filepath,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.CreatedBy,
		&i.Username,
		&i.Content,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Filepath,
	)
	return i, err
}

const getOne = `-- name: GetOne :one
select posts.id, created_by, username, posts.content, posts.created_at, posts.updated_at, filepath, comments.id, post_id, user_id, comments.content, comments.created_at, comments.updated_at from posts JOIN comments on posts.id = comments.post_id where posts.id = $1
`

type GetOneRow struct {
	ID          string
	CreatedBy   string
	Username    string
	Content     string
	CreatedAt   pgtype.Timestamp
	UpdatedAt   pgtype.Timestamp
	Filepath    pgtype.Text
	ID_2        string
	PostID      string
	UserID      string
	Content_2   string
	CreatedAt_2 pgtype.Timestamp
	UpdatedAt_2 pgtype.Timestamp
}

func (q *Queries) GetOne(ctx context.Context, id string) (GetOneRow, error) {
	row := q.db.QueryRow(ctx, getOne, id)
	var i GetOneRow
	err := row.Scan(
		&i.ID,
		&i.CreatedBy,
		&i.Username,
		&i.Content,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Filepath,
		&i.ID_2,
		&i.PostID,
		&i.UserID,
		&i.Content_2,
		&i.CreatedAt_2,
		&i.UpdatedAt_2,
	)
	return i, err
}

const getPostsByUser = `-- name: GetPostsByUser :many
SELECT id, created_by, username, content, created_at, updated_at, filepath FROM posts WHERE created_by = $1 OR username = $2 ORDER BY created_at DESC LIMIT $3 OFFSET $4
`

type GetPostsByUserParams struct {
	CreatedBy string
	Username  string
	Limit     int32
	Offset    int32
}

func (q *Queries) GetPostsByUser(ctx context.Context, arg GetPostsByUserParams) ([]Post, error) {
	rows, err := q.db.Query(ctx, getPostsByUser,
		arg.CreatedBy,
		arg.Username,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.CreatedBy,
			&i.Username,
			&i.Content,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Filepath,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
