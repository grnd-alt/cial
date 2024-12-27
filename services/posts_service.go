package services

import (
	"backendsetup/m/db/sql/dbgen"
	"context"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type PostsService struct {
	query       *dbgen.Queries
	fileService *FileService
}

func InitPostsService(queries *dbgen.Queries, fileFileService *FileService) *PostsService {
	return &PostsService{
		query:       queries,
		fileService: fileFileService,
	}
}

func (n *PostsService) CreatePost(username string, createdBy string, content string, file *multipart.FileHeader) (*dbgen.Post, error) {
	fileReader, err := file.Open()
	if err != nil {
		return nil, err
	}
	id := uuid.NewString()
	location, err := n.fileService.UploadFile(id, fileReader, file.Size, file.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	post, err := n.query.CreatePost(context.Background(), dbgen.CreatePostParams{
		ID:        id,
		CreatedBy: createdBy,
		Content:   content,
		Username:  username,
		Filepath:  pgtype.Text{String: location, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return &post, err
}

func (n *PostsService) GetPosts(createdBy string) ([]dbgen.Post, error) {
	posts, err := n.query.GetPostsByUser(context.Background(), createdBy)
	return posts, err
}
