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

func (c *CommentsService) GetComments(postID string, page int32) ([]dbgen.Comment, error) {
	comments, err := c.Queries.GetCommentsByPost(context.Background(), dbgen.GetCommentsByPostParams{Offset: 10 * page, PostID: postID, Limit: 10})
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (c *CommentsService) CreateComment(userId string, username string, content string, postID string) (*dbgen.Comment, error) {
	comment, err := c.Queries.CreateComment(context.Background(), dbgen.CreateCommentParams{
		ID:       uuid.NewString(),
		UserID:   userId,
		Content:  content,
		PostID:   postID,
		UserName: username,
	})
	return &comment, err
}
