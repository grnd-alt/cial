package routes

import (
	"backendsetup/m/config"
	"backendsetup/m/controllers"
	"backendsetup/m/db/sql/dbgen"
	"backendsetup/m/middleware"
	"backendsetup/m/services"
	"fmt"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Init(verifier *oidc.IDTokenVerifier, conf *config.Config, queries *dbgen.Queries) *gin.Engine {
	if conf.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	corsConf := cors.DefaultConfig()
	corsConf.AddAllowHeaders("Authorization")
	corsConf.AllowCredentials = true
	corsConf.AllowAllOrigins = true
	fileService, err := services.InitFileService(conf.S3URL, conf.S3AccessKey, conf.S3SecretKey, conf.S3BucketName, conf.AppEnv)
	for {
		if err != nil {
			fmt.Printf("s3 init failed %v\n", err)
			time.Sleep(5 * time.Second)
			fileService, err = services.InitFileService(conf.S3URL, conf.S3AccessKey, conf.S3SecretKey, conf.S3BucketName, conf.AppEnv)
			continue
		}
		break
	}
	postService := services.InitPostsService(queries, fileService)
	commentsService := services.InitCommentsService(queries)

	userController := controllers.InitUserController(conf, verifier)
	postsController := controllers.InitPostsController(conf, postService)
	commentsController := controllers.InitCommentsController(commentsService)

	engine.Use(cors.New(corsConf))
	engine.Use(gin.Logger())

	engine.GET("/api/hello", controllers.HelloWorldHandler)

	engine.Use(middleware.ProtectedMiddleware(verifier, conf.AppEnv))
	{
		engine.POST("/api/comments/create", commentsController.CreateComment)
		engine.GET("/api/me", userController.Me)
		engine.POST("/api/posts/create", postsController.Create)
		engine.DELETE("/api/posts/:postId", postsController.Delete)
		engine.GET("/api/posts", postsController.GetLatest)
		engine.GET("/api/posts/:username", postsController.GetPostsByUser)
	}
	return engine
}
