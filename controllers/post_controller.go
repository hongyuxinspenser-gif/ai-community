package controllers

import (
	"ai-community/config"
	"ai-community/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreatePost 创建帖子
func CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&post)
	c.JSON(http.StatusOK, post)
}

// GetPosts 获取帖子列表
func GetPosts(c *gin.Context) {
	var posts []models.Post
	config.DB.Find(&posts)
	c.JSON(http.StatusOK, gin.H{"data": posts})
}

// DeletePost 删除帖子
func DeletePost(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Post{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}