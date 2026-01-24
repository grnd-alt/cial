package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"sync"

	"backendsetup/m/db/sql/dbgen"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type PostsService struct {
	query               *dbgen.Queries
	fileService         *FileService
	notificationService *NotificationService
}

func InitPostsService(queries *dbgen.Queries, fileService *FileService, notificationService *NotificationService) *PostsService {
	return &PostsService{
		query:               queries,
		fileService:         fileService,
		notificationService: notificationService,
	}
}

func (n *PostsService) GetPost(postId string) (*PostWithComments, error) {
	post, err := n.query.GetOne(context.Background(), postId)
	if err != nil {
		return nil, err
	}
	fullpost, err := n.populatePosts([]dbgen.Post{post})
	if err != nil {
		return nil, err
	}
	return &fullpost[0], nil
}

func (n *PostsService) CreatePost(username string, createdBy string, content string, file *multipart.FileHeader) (*PostWithComments, error) {
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

	posts, err := n.populatePosts([]dbgen.Post{post})
	if err != nil {
		return nil, err
	}
	data := NotificationData{
		Type:  NewPostNotificationType,
		Title: "There's a new post!",
		Body:  fmt.Sprintf("%s just posted", username),
		Data: NewPostData{
			Author: username,
		},
	}
	go n.notificationService.SendFollowersNotification(data, createdBy)
	return &posts[0], err
}

func (n *PostsService) DeletePost(postId string) error {
	err := n.query.DeletePost(context.Background(), postId)
	if err != nil {
		return err
	}
	return n.query.DeleteCommentsByPost(context.Background(), postId)
}

type PostWithComments struct {
	Post     *dbgen.Post
	Comments []*dbgen.Comment
}

func (n *PostsService) populatePosts(posts []dbgen.Post) ([]PostWithComments, error) {
	ids := make([]string, len(posts))
	for i, post := range posts {
		ids[i] = post.ID
	}
	comments, err := n.query.GetCommentsByPosts(context.Background(), ids)
	if err != nil {
		return nil, err
	}
	commentsMap := make(map[string][]*dbgen.Comment)
	for _, comment := range comments {
		commentsMap[comment.PostID] = append(commentsMap[comment.PostID], &comment)
	}
	postsWithComments := make([]PostWithComments, len(posts))

	var wg sync.WaitGroup
	for i, post := range posts {
		post := post
		wg.Add(1)
		go func(post dbgen.Post) {
			defer wg.Done()
			filepath, err := n.fileService.GetFileUrl(post.ID)
			if err != nil {
				return
			}
			post.Filepath = pgtype.Text{String: filepath, Valid: true}
			postsWithComments[i] = PostWithComments{
				Post:     &post,
				Comments: commentsMap[post.ID],
			}
		}(post)
	}
	wg.Wait()
	return postsWithComments, nil
}

func (n *PostsService) GetFeed(offset int32) ([]PostWithComments, error) {
	posts, err := n.query.GetLatestPosts(context.Background(), dbgen.GetLatestPostsParams{Limit: 10, Offset: offset * 10})
	if err != nil {
		return nil, err
	}
	return n.populatePosts(posts)
}

func (n *PostsService) GetPosts(createdBy string, page int, username string) ([]PostWithComments, error) {
	posts, err := n.query.GetPostsByUser(context.Background(), dbgen.GetPostsByUserParams{Username: username, Limit: 10, Offset: int32(page * 10)})
	if err != nil {
		return nil, err
	}

	return n.populatePosts(posts)
}
