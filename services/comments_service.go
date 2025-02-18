package services

import (
	"backendsetup/m/db/sql/dbgen"
	"context"

	"github.com/google/uuid"
)

type CommentsService struct {
	Queries *dbgen.Queries
}

func InitCommentsService(queries *dbgen.Queries) *CommentsService {
	return &CommentsService{
		Queries: queries,
	}
}

func (c *CommentsService) CreateComment(userId string, username string,content string, postID string) (*dbgen.Comment, error) {
	comment, err := c.Queries.CreateComment(context.Background(), dbgen.CreateCommentParams{
		ID: uuid.NewString(),
		UserID: userId,
		Content:   content,
		PostID:    postID,
		UserName:  username,
	})
	return &comment, err
}
