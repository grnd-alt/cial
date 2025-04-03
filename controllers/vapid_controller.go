package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type VapidController struct {
	pubkey string
}

func InitVapidController(pubkey string) *VapidController {
	return &VapidController{pubkey: pubkey}
}

func (v *VapidController) GetPublicKey(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"publicKey": v.pubkey})
}
