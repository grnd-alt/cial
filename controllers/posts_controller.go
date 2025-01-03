package controllers

import (
	"backendsetup/m/config"
	"backendsetup/m/middleware"
	"backendsetup/m/services"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreatePostJSON struct {
	Content string `form:"content" binding:"required"`
}

type PostsController struct {
	PostsService *services.PostsService
}

func InitPostsController(conf *config.Config, postsService *services.PostsService) *PostsController {
	return &PostsController{
		PostsService: postsService,
	}
}

func (n *PostsController) Create(ctx *gin.Context) {
	claimsInterface, exists := ctx.Get("claims")
	if !exists {
		response := gin.H{"error": "claims not found"}
		fmt.Println(response)
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}
	claims, ok := claimsInterface.(middleware.Claims)
	if !ok {
		response := gin.H{"error": "invalid claims"}
		fmt.Println(response)
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	var params CreatePostJSON
	if err := ctx.ShouldBind(&params); err != nil {
		response := gin.H{"error": err.Error()}
		fmt.Println(response)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		response := gin.H{"error": err.Error()}
		fmt.Println(response)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	post, err := n.PostsService.CreatePost(claims.Username, claims.Sub, params.Content, fileHeader)
	if err != nil {
		response := gin.H{"error": err.Error()}
		fmt.Println(response)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	ctx.JSON(http.StatusOK, post)
}

func (n *PostsController) GetPostsByUser(ctx *gin.Context) {
	claimsInterface, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "claims not found"})
		return
	}
	claims, ok := claimsInterface.(middleware.Claims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
		return
	}

	username := ctx.Param("username")
	pageStr := ctx.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 0 {
		page = 0
	}
	posts, err := n.PostsService.GetPosts(claims.Sub, page, username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, posts)
}
