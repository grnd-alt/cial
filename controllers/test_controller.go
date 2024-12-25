package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HelloWorldHandler(ctx *gin.Context){
	ctx.String(http.StatusOK, "hello world")
}
