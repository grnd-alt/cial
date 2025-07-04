package controllers

import (
	"backendsetup/m/middleware"
	"backendsetup/m/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentsController struct {
	CommentsService *services.CommentsService
}

func InitCommentsController(commentsService *services.CommentsService) *CommentsController {
	return &CommentsController{
		CommentsService: commentsService,
	}
}

func (c *CommentsController) GetCommentsByPost(ctx *gin.Context) {
	postID := ctx.Param("postId")
	page, err := strconv.Atoi(ctx.Query("page"))
	if postID == "" || err != nil || page < 0 {
		response := gin.H{"error": "incorrect parameters"}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	comments, err := c.CommentsService.GetComments(postID, int32(page))
	if err != nil {
		response := gin.H{"error": err.Error()}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"comments": comments,
	})
}

func (c *CommentsController) CreateComment(ctx *gin.Context) {
	claimsInterface, exists := ctx.Get("claims")
	if !exists {
		response := gin.H{"error": "claims not found"}
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}
	claims, ok := claimsInterface.(middleware.Claims)
	if !ok {
		response := gin.H{"error": "invalid claims"}
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	var requestBody struct {
		Content string `json:"content"`
		PostID  string `json:"post_id"`
	}
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		response := gin.H{"error": "invalid request body"}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	comment, err := c.CommentsService.CreateComment(claims.Sub, claims.Username, requestBody.Content, requestBody.PostID)
	if err != nil {
		response := gin.H{"error": err.Error()}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusOK, comment)
}
