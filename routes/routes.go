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
	engine := gin.New()
	corsConf := cors.DefaultConfig()
	corsConf.AddAllowHeaders("Authorization")
	corsConf.AllowCredentials = true
	corsConf.AllowAllOrigins = true

	fileService, err := services.InitFileService(conf.S3URL, conf.S3AccessKey, conf.S3SecretKey, conf.S3BucketName)
	for {
		if err != nil {
			fmt.Printf("s3 init failed %v", err)
			time.Sleep(5 * time.Second)
			fileService, err = services.InitFileService(conf.S3URL, conf.S3AccessKey, conf.S3SecretKey, conf.S3BucketName)
			continue
		}
		break
	}
	postService := services.InitPostsService(queries, fileService)

	userController := controllers.InitUserController(conf, verifier)
	postsController := controllers.InitPostsController(conf, postService)

	engine.Use(cors.New(corsConf))
	engine.Use(gin.Logger())

	engine.GET("/api/hello", controllers.HelloWorldHandler)

	engine.Use(middleware.ProtectedMiddleware(verifier, conf.AppEnv))
	{
		engine.GET("/api/me", userController.Me)
		engine.POST("/api/posts/create", postsController.Create)
		engine.GET("/api/posts", postsController.GetPostsByUser)
	}
	return engine
}
