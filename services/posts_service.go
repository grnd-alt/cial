package services

import (
	"backendsetup/m/db/sql/dbgen"
	"context"
	"mime/multipart"
	"sync"

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

type PostWithComments struct {
	Post     *dbgen.Post
	Comments []*dbgen.GetCommentsByPostsRow
}

func (n *PostsService) populatePosts(posts []dbgen.Post) ([]PostWithComments,error) {
	ids := make([]string, len(posts))
	for i, post := range posts {
		ids[i] = post.ID
	}
	comments, err := n.query.GetCommentsByPosts(context.Background(), ids)
	if err != nil {
		return nil, err
	}
	commentsMap := make(map[string][]*dbgen.GetCommentsByPostsRow)
	for _, comment := range comments {
		commentsMap[comment.PostID] = append(commentsMap[comment.PostID], &comment)
	}
	postsWithComments := make([]PostWithComments, len(posts))

	var mut sync.Mutex
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
			mut.Lock()
			postsWithComments[i] = PostWithComments{
				Post:     &post,
				Comments: commentsMap[post.ID],
			}
			mut.Unlock()
		}(post)
	}
	wg.Wait()
	return postsWithComments,nil
}


func (n *PostsService) GetFeed(offset int32) ([]PostWithComments, error) {
	posts, err := n.query.GetLatestPosts(context.Background(), dbgen.GetLatestPostsParams{Limit: 10, Offset: offset})
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
