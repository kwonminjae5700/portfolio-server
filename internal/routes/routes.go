package routes

import (
	"portfolio-server/internal/handlers"
	"portfolio-server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "Server is running",
		})
	})

	authHandler := handlers.NewAuthHandler()
	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/profile", middleware.AuthMiddleware(), authHandler.GetProfile)
	}

	articleHandler := handlers.NewArticleHandler()
	articles := router.Group("/articles")
	{
		articles.GET("", articleHandler.GetArticles)
		articles.GET("/:id", articleHandler.GetArticle)
		articles.POST("", middleware.AuthMiddleware(), articleHandler.CreateArticle)
		articles.PUT("/:id", middleware.AuthMiddleware(), articleHandler.UpdateArticle)
		articles.DELETE("/:id", middleware.AuthMiddleware(), articleHandler.DeleteArticle)

		commentHandler := handlers.NewCommentHandler()
		articles.GET("/:article_id/comments", commentHandler.GetCommentsByArticle)
	}

	commentHandler := handlers.NewCommentHandler()
	comments := router.Group("/comments")
	{
		comments.POST("", middleware.AuthMiddleware(), commentHandler.CreateComment)
		comments.PUT("/:id", middleware.AuthMiddleware(), commentHandler.UpdateComment)
		comments.DELETE("/:id", middleware.AuthMiddleware(), commentHandler.DeleteComment)
	}
}
