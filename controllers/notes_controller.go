package controllers

import (
	"backendsetup/m/config"
	"backendsetup/m/middleware"
	"backendsetup/m/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateNoteJSON struct {
	Title     string `json:"Title"`
	Content   string `json:"Content"`
}

type NotesController struct {
	NotesService *services.NotesService
}

func InitNotesController(conf *config.Config, notesService *services.NotesService) *NotesController {
	return &NotesController{
		NotesService: notesService,
	}
}

func (n *NotesController) Create(ctx *gin.Context) {
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
	var params CreateNoteJSON
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	note, err := n.NotesService.CreateNote(claims.Sub, params.Title, params.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, note)
}

func (n *NotesController) GetNotesByUser(ctx *gin.Context) {
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
	notes, err := n.NotesService.GetNotes(claims.Sub)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, notes)
}
