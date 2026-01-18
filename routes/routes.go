package routes

import (
	"ai-community/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		// Post 路由
		api.POST("/posts", controllers.CreatePost)
		api.GET("/posts", controllers.GetPosts)
		api.DELETE("/posts/:id", controllers.DeletePost)

		// Comment 路由
		api.POST("/comments", controllers.CreateComment)
		api.DELETE("/comments/:id", controllers.DeleteComment)
		api.GET("/posts/:post_id/comments", controllers.GetCommentsByPost)
	}

	return r
}