// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package dbgen

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Comment struct {
	ID        string
	PostID    string
	UserID    string
	Content   string
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
	UserName  string
}

type Post struct {
	ID        string
	CreatedBy string
	Username  string
	Content   string
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
	Filepath  pgtype.Text
}
