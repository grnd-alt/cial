package routes

import (
	"backendsetup/m/config"
	"backendsetup/m/controllers"
	"backendsetup/m/db/sql/dbgen"
	"backendsetup/m/middleware"
	"backendsetup/m/services"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Init(verifier *oidc.IDTokenVerifier, conf *config.Config,queries *dbgen.Queries) *gin.Engine {
	engine := gin.New()
	corsConf:=cors.DefaultConfig()
	corsConf.AddAllowHeaders("Authorization")
	corsConf.AllowCredentials=true
	corsConf.AllowAllOrigins=true

	notesService := services.InitNotesService(queries)

	userController := controllers.InitUserController(conf, verifier)
	notesController := controllers.InitNotesController(conf, notesService)

	engine.Use(cors.New(corsConf))
	engine.Use(gin.Logger())

	engine.GET("/api/hello", controllers.HelloWorldHandler)

	engine.Use(middleware.ProtectedMiddleware(verifier))
	{
		engine.GET("/api/me", userController.Me)
		engine.POST("/api/notes/create", notesController.Create)
		engine.GET("/api/notes", notesController.GetNotesByUser)
	}
	return engine
}
