package controllers

import (
	"backendsetup/m/config"
	"backendsetup/m/middleware"
	"backendsetup/m/services"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	config              *config.Config
	verifier            *oidc.IDTokenVerifier
	followService       *services.FollowService
	userService         *services.UserService
	notificationService *services.NotificationService
}

func InitUserController(config *config.Config, verifier *oidc.IDTokenVerifier, followService *services.FollowService, userService *services.UserService, notificationService *services.NotificationService) *UserController {
	return &UserController{
		config:              config,
		verifier:            verifier,
		followService:       followService,
		userService:         userService,
		notificationService: notificationService,
	}
}

func getBearer(header string) (string, error) {
	vals := strings.Split(header, " ")
	if len(vals) != 2 {
		return "", errors.New("Invalid Authorization header")
	}
	if vals[0] != "Bearer" {
		return "", errors.New("Invalid Authorization header")
	}
	return vals[1], nil
}

func (u UserController) Follow(ctx *gin.Context) {
	type body = struct {
		Subscription any `json:"subscription"`
	}
	var reqBody body
	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "claims not found"})
		return
	}
	username := claims.(middleware.Claims).Username
	followingname := ctx.Param("username")
	if err := ctx.BindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	subscriptionData, err := json.Marshal(reqBody.Subscription)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = u.followService.Follow(username, followingname, subscriptionData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Followed"})
}

func (u UserController) Unfollow(ctx *gin.Context) {
	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "claims not found"})
		return
	}
	username := claims.(middleware.Claims).Username
	followingname := ctx.Param("username")
	err := u.followService.Unfollow(username, followingname)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Unfollowed"})
}

func (u UserController) GetUser(ctx *gin.Context) {
	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "claims not found"})
		return
	}
	type userdata = struct {
		FollowingCount int64 `json:"followingcount"`
		Followers      int64 `json:"followers"`
		IsFollowing    bool  `json:"isfollowing"`

		Username string `json:"username"`
	}
	username := ctx.Param("username")
	followingCount, err := u.followService.GetFollowingCount(username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	followerCount, err := u.followService.GetFollowersCount(username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	isfollowing := u.followService.IsFollowing(claims.(middleware.Claims).Username, ctx.Param("username"))

	ctx.JSON(http.StatusOK, userdata{Followers: followerCount, FollowingCount: followingCount, Username: username, IsFollowing: isfollowing})
}

func (u UserController) GetFollowers(ctx *gin.Context) {
	username := ctx.Param("username")
	followers, err := u.followService.GetFollowers(username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, followers)
}

func (u *UserController) UpdateBrowserData(ctx *gin.Context) {
	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "claims not found"})
		return
	}

	type body = struct {
		Subscription json.RawMessage `json:"subscription"`
	}
	var reqBody body
	err := ctx.BindJSON(&reqBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pclaims := claims.(middleware.Claims)

	err = u.userService.InsertSubscription(pclaims.Sub, reqBody.Subscription)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Subscription added"})
	return
}

func (u *UserController) NotifyMe(ctx *gin.Context) {
	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "claims not found"})
		return
	}
	userId := claims.(middleware.Claims).Sub
	err := u.notificationService.SendNotification("You have a new notification", userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Notification sent"})
}

func (u *UserController) Me(ctx *gin.Context) {
	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "claims not found"})
		return
	}
	u.userService.CreateUserIfNotExists(claims.(middleware.Claims).Username, claims.(middleware.Claims).Sub)
	ctx.JSON(http.StatusOK, claims)
}
