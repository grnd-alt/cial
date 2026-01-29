package routes

import (
	"backendsetup/m/config"
	"backendsetup/m/controllers"
	"backendsetup/m/db/sql/dbgen"
	"backendsetup/m/middleware"
	"backendsetup/m/services"
	"backendsetup/m/workers"
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

	fileURLCache := services.InitInMemoryUrlCache()
	fileService, err := services.InitFileService(conf.S3URL, conf.S3AccessKey, conf.S3SecretKey, conf.S3BucketName, fileURLCache, conf.AppEnv)
	for {
		if err != nil {
			fmt.Printf("s3 init failed %v\n", err)
			time.Sleep(5 * time.Second)
			fileService, err = services.InitFileService(conf.S3URL, conf.S3AccessKey, conf.S3SecretKey, conf.S3BucketName, fileURLCache, conf.AppEnv)
			continue
		}
		break
	}

	notificationService := services.InitNotificationServe(conf, queries)
	followService := services.InitFollowService(queries)
	userService := services.InitUserService(queries)
	postService := services.InitPostsService(queries, fileService, notificationService)
	commentsService := services.InitCommentsService(queries)
	countersService := services.InitCountersService(queries)

	notificationWorker := workers.NewReminderNotificationWorker(queries,notificationService)
	go notificationWorker.StartWorker()

	userController := controllers.InitUserController(conf, verifier, followService, userService, notificationService)
	postsController := controllers.InitPostsController(conf, postService)
	commentsController := controllers.InitCommentsController(commentsService)
	vapidController := controllers.InitVapidController(conf.VAPIDPub)
	countersController := controllers.InitCountersController(countersService)

	engine.Use(cors.New(corsConf))
	engine.Use(gin.Logger())

	engine.GET("/api/hello", controllers.HelloWorldHandler)
	engine.GET("/api/vapid", vapidController.GetPublicKey)

	engine.Use(middleware.ProtectedMiddleware(verifier, conf.AppEnv, queries))
	{

		// posts
		engine.POST("/api/posts/create", postsController.Create)
		engine.DELETE("/api/posts/:postId", postsController.Delete)
		engine.GET("/api/posts", postsController.GetLatest)
		engine.GET("/api/posts/:username", postsController.GetPostsByUser)
		// comments
		engine.POST("/api/comments/create", commentsController.CreateComment)
		engine.GET("/api/comments/:postId", commentsController.GetCommentsByPost)

		// users
		engine.GET("/api/me", userController.Me)
		engine.GET("/api/user/find", userController.FindUser)
		engine.GET("/api/users/:username", userController.GetUser)
		engine.GET("/api/users/:username/followers", userController.GetFollowers)
		engine.POST("/api/users/update-browser-data", userController.UpdateBrowserData)
		engine.POST("/api/users/follow/:username", userController.Follow)
		engine.POST("/api/users/unfollow/:username", userController.Unfollow)
		engine.POST("/api/users/notifyme", userController.NotifyMe)

		// counters
		engine.POST("api/counters", countersController.CreateCounter)
		engine.POST("api/counters/:counterID/events/add", countersController.AddEvent)
		engine.GET("api/counters/:counterID/events", countersController.GetEvents)
		engine.GET("api/counters", countersController.GetCounters)
		engine.GET("api/counters/:counterID", countersController.GetCounter)
		engine.POST("api/counters/share", countersController.ShareCounter)
	}
	return engine
}
