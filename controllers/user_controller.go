package controllers

import (
	"backendsetup/m/config"
	"errors"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	config   *config.Config
	verifier *oidc.IDTokenVerifier
}

func InitUserController(config *config.Config, verifier *oidc.IDTokenVerifier) *UserController {
	return &UserController{
		config:   config,
		verifier: verifier,
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

func (u UserController) Me(ctx *gin.Context) {
	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "claims not found"})
		return
	}
	ctx.JSON(http.StatusOK, claims)
}
