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
	commentHandler := handlers.NewCommentHandler()
	articles := router.Group("/articles")
	{
		articles.GET("", articleHandler.GetArticles)
		articles.GET("/:id", articleHandler.GetArticle)
		articles.POST("", middleware.AuthMiddleware(), articleHandler.CreateArticle)
		articles.PUT("/:id", middleware.AuthMiddleware(), articleHandler.UpdateArticle)
		articles.DELETE("/:id", middleware.AuthMiddleware(), articleHandler.DeleteArticle)
		
		// Comments routes
		articles.GET("/:id/comments", commentHandler.GetComments)
		articles.POST("/:id/comments", middleware.AuthMiddleware(), commentHandler.CreateComment)
		articles.PUT("/:id/comments/:commentId", middleware.AuthMiddleware(), commentHandler.UpdateComment)
		articles.DELETE("/:id/comments/:commentId", middleware.AuthMiddleware(), commentHandler.DeleteComment)
	}

	categoryHandler := handlers.NewCategoryHandler()
	categories := router.Group("/categories")
	{
		categories.GET("", categoryHandler.GetCategories)
		categories.GET("/:id", categoryHandler.GetCategory)
		categories.POST("", middleware.AuthMiddleware(), categoryHandler.CreateCategory)
		categories.PUT("/:id", middleware.AuthMiddleware(), categoryHandler.UpdateCategory)
		categories.DELETE("/:id", middleware.AuthMiddleware(), categoryHandler.DeleteCategory)
	}
}
